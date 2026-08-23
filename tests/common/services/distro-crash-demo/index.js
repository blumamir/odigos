const http = require('http');

const PORT = 3000;
const CRASH_DELAY_ENV = 'CRASH_DELAY_SECONDS';
let requestCount = 0;

function getCrashDelaySeconds() {
  const rawValue = process.env[CRASH_DELAY_ENV];
  if (rawValue === undefined || rawValue === '') {
    return 0;
  }

  const parsed = Number.parseInt(rawValue, 10);
  if (Number.isNaN(parsed) || parsed < 0) {
    console.log(`⚠️  Invalid ${CRASH_DELAY_ENV}=${rawValue}, defaulting to 0`);
    return 0;
  }

  return parsed;
}

function isInstrumented() {
  return Boolean(process.env.ODIGOS_DISTRO_NAME);
}

function crashOnInstrumentation(delaySeconds) {
  console.log('💥 ODIGOS INSTRUMENTATION DETECTED - CRASHING!');
  console.log(`   ODIGOS_DISTRO_NAME=${process.env.ODIGOS_DISTRO_NAME}`);
  if (delaySeconds > 0) {
    console.log(`   Crash delay: ${delaySeconds}s (${CRASH_DELAY_ENV}=${process.env[CRASH_DELAY_ENV]})`);
  }

  console.log('💀 This application is incompatible with Odigos instrumentation');
  console.log('🔄 Expecting Odigos auto-rollback to uninstrument this service...');

  process.exit(1);
}

function scheduleCrashOnInstrumentation(delaySeconds) {
  if (delaySeconds === 0) {
    crashOnInstrumentation(delaySeconds);
    return;
  }

  console.log(`⏳ Instrumentation detected - will crash in ${delaySeconds}s`);
  setTimeout(() => crashOnInstrumentation(delaySeconds), delaySeconds * 1000);
}

function writeProbeResponse(res, endpoint) {
  res.writeHead(200, { 'Content-Type': 'text/plain' });
  res.end(`${endpoint} ok\n`);
}

const server = http.createServer((req, res) => {
  if (req.method === 'GET' && req.url === '/healthz') {
    writeProbeResponse(res, 'healthz');
    return;
  }

  if (req.method === 'GET' && req.url === '/readyz') {
    writeProbeResponse(res, 'readyz');
    return;
  }

  requestCount++;

  const response = {
    message: 'Hello from distro crash demo service!',
    requestCount,
    instrumented: isInstrumented(),
    distroName: process.env.ODIGOS_DISTRO_NAME || null,
    crashDelaySeconds: isInstrumented() ? getCrashDelaySeconds() : null,
    timestamp: new Date().toISOString(),
    status: isInstrumented()
      ? 'running before scheduled crash'
      : 'healthy - no instrumentation detected',
  };

  res.writeHead(200, {
    'Content-Type': 'application/json',
    'X-Request-Count': requestCount.toString(),
    'X-Instrumented': isInstrumented().toString(),
  });
  res.end(JSON.stringify(response, null, 2) + '\n');

  console.log(`📝 Request ${requestCount}: ${req.method} ${req.url}`);
});

process.on('SIGTERM', () => {
  console.log('🛑 SIGTERM received, shutting down gracefully...');
  server.close(() => {
    console.log('✅ Server closed');
    process.exit(0);
  });
});

process.on('SIGINT', () => {
  console.log('🛑 SIGINT received, shutting down gracefully...');
  server.close(() => {
    console.log('✅ Server closed');
    process.exit(0);
  });
});

server.listen(PORT, () => {
  console.log(`🚀 Distro Crash Demo Service started on port ${PORT}`);
  console.log(`📊 Process ID: ${process.pid}`);
  console.log('🔍 Checking for Odigos instrumentation at startup...');

  if (isInstrumented()) {
    console.log('⚠️  Odigos instrumentation detected at startup!');
    scheduleCrashOnInstrumentation(getCrashDelaySeconds());
  } else {
    console.log('✅ No Odigos instrumentation detected - service running normally');
  }
});

process.on('uncaughtException', (err) => {
  console.error('💀 Uncaught Exception:', err);
  process.exit(1);
});

process.on('unhandledRejection', (reason, promise) => {
  console.error('💀 Unhandled Rejection at:', promise, 'reason:', reason);
  process.exit(1);
});
