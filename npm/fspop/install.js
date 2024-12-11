const os = require('os')
const path = require('path')
const fs = require('fs')
const https = require("https");
const zlib = require("zlib");
const packageVersion = require(path.join(__dirname, "package.json")).version;

function getPlatformPackage() {
	const platform = os.platform()
	const arch = os.arch()

	const platforms = {
		'win32': {
			'x64': 'win32-x64',
			'arm64': 'win32-arm64'
		},
		'darwin': {
			'x64': 'darwin-x64',
			'arm64': 'darwin-arm64'
		},
		'linux': {
			'x64': 'linux-x64',
			'arm64': 'linux-arm64'
		}
	}

	const platformSupport = platforms[platform]
	if (!platformSupport) throw new Error(`Unsupported platform: ${platform}`)

	const binary = platformSupport[arch]
	if (!binary) throw new Error(`Unsupported architecture: ${arch}`)

	return `fspop-${binary}`
}

function fetch(url) {
	return new Promise((resolve, reject) => {
		https.get(url, (res) => {
			if ((res.statusCode === 301 || res.statusCode === 302) && res.headers.location)
				return fetch(res.headers.location).then(resolve, reject);

			if (res.statusCode !== 200)
				return reject(new Error(`Server responded with ${res.statusCode}`));

			let chunks = [];
			res.on("data", (chunk) => chunks.push(chunk));
			res.on("end", () => resolve(Buffer.concat(chunks)));
		}).on("error", reject);
	});
}

function extractFileFromTarGzip(buffer, subpath) {
	try {
		buffer = zlib.unzipSync(buffer);
	} catch (err) {
		throw new Error(`[fspop] Invalid gzip data in archive: ${err && err.message || err}`);
	}

	let str = (i, n) => String.fromCharCode(...buffer.subarray(i, i + n)).replace(/\0.*$/, "");
	let offset = 0;
	subpath = `package/${subpath}`;
	while (offset < buffer.length) {
		let name = str(offset, 100);
		let size = parseInt(str(offset + 124, 12), 8);
		offset += 512;
		if (!isNaN(size)) {
			if (name === subpath) return buffer.subarray(offset, offset + size);
			offset += size + 511 & ~511;
		}
	}

	throw new Error(`[fspop] Could not find ${JSON.stringify(subpath)} in archive`);
}

async function downloadDirectlyFromNPM(packageName, subpath, binPath) {
	const url = `https://registry.npmjs.org/${packageName}/-/${packageName}-${packageVersion}.tgz`;
	console.debug(`[fspop] Trying to download ${JSON.stringify(url)}`);
	try {
		fs.writeFileSync(binPath, extractFileFromTarGzip(await fetch(url), subpath));
		fs.chmodSync(binPath, 493);
	} catch (e) {
		console.error(`[fspop] Failed to download ${JSON.stringify(url)}: ${e && e.message || e}`);
		throw e;
	}
}

async function install() {
	try {
		const platformPackage = getPlatformPackage()
		const binaryFileName = `fspop${(os.platform() === 'win32' ? '.exe' : '')}`

		downloadDirectlyFromNPM(
			platformPackage,
			binaryFileName,
			path.join(__dirname, `bin/_${binaryFileName}`)
		)

		console.log('[fspop] Binary installed successfully')
	} catch (error) {
		console.error('[fspop] Failed to install binary:', error)
		process.exit(1)
	}
}

install().catch(err => {
	console.error(err)
	process.exit(1)
})
