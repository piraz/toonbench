import Benchmark from "benchmark";
import { encode, decode } from "@toon-format/toon";

if (typeof global.gc !== 'function') {
    console.error('Run node with --expose-gc to get reliable memory numbers: node --expose-gc src/bench.js');
    process.exit(1);
}

const NUM_USERS = 100_000;

const roles = ["admin", "user", "moderator", "superuser"]

function generatePayload(n) {
    const users = [];
    for (let i= 0; i < n; i++ ) {
        users.push({
            id: i + 1,
            name: ``,
            role: roles[i % 4]
        });
    }
    return { users };
}


const payload  = generatePayload(NUM_USERS);
const jsonStr = JSON.stringify(payload);
const toonStr = encode(payload);

function formatMB(bytes) {
    return (bytes / 1024 / 1024).toFixed(3);
}

function nsPerOpFromHz(hz) {
    if (!hz || hz <= 0) return NaN;
    return 1e9 / hz;
}

function runSingleBenchmark(name, fn) {
    return new Promise((resolve) => {
        // force GC and take baseline
        global.gc();
        const memBefore = process.memoryUsage().heapUsed;
        const suite = new Benchmark.Suite();

        suite
            .add(name, fn)
            .on("cycle", (event) => {
                const b = event.target;
            })
            .on("complete", function() {
                global.gc();
                // TODO: fix memory utilization
                const memAfter = process.memoryUsage().heapUsed;
                const memDelta = Math.max(0, memAfter - memBefore);

                const bench = this[0];
                // count: number of times executed (fallback to estimation)
                const executed = bench.count || Math.round(bench.hz * bench.times.elapsed) || 0;
                const hz = bench.hz || 0;
                const nsOp = nsPerOpFromHz(hz);
                // TODO: fix memory utilization
                const bPerOp = executed > 0 ? (memDelta / executed) : NaN;

                console.log(`\nRESULT: ${name}`);
                console.log(`  ops/sec: ${hz.toFixed(2)}`);
                if (!Number.isNaN(nsOp)) console.log(`  ns/op: ${Math.round(nsOp).toLocaleString()}`);
                console.log(`  executed (approx): ${executed}`);
                // TODO: fix memory utilization
                // console.log(`  mem delta: ${memDelta} bytes (${formatMB(memDelta)} MB)`);
                // console.log(`  approx B/op: ${Number.isNaN(bPerOp) ? 'N/A' : Math.round(bPerOp) + ' bytes'}`);

                resolve({
                    name,
                    hz,
                    nsOp,
                    executed,
                    memDelta,
                    bPerOp,
                });
            })
            .run({ async: false });

    });
}

async function main() {
    console.log('Starting benchmarks (Node.js + benchmark.js). Ensure minimal other activity for best results.');
    console.log('Payload users =', NUM_USERS);

    // Define the functions to benchmark
    const tests = [
        {
            name: 'JSON.stringify',
            fn: () => JSON.stringify(payload)
        },
        {
            name: 'JSON.parse',
            fn: () => JSON.parse(jsonStr)
        },
        {
            name: 'TOON.encode',
            fn: () => encode(payload)
        },
        {
            name: 'TOON.decode',
            fn: () => decode(toonStr)
        },
    ];

    const results = [];
    for (const t of tests) {
        // Run each benchmark isolated to better measure memory impact
        console.log(`\n--- Running: ${t.name} ---`);
        const r = await runSingleBenchmark(t.name, t.fn);
        results.push(r);
    }

    console.log('\nAll done.');
    // You can further process 'results' to build a table if desired.
}

main().catch(err => {
    console.error("Error:", err);
    process.exit(1);
});
