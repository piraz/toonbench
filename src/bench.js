import Benchmark from "benchmark";
import { encode, decode } from "@toon-format/toon";

import { createRequire } from "module";
const require = createRequire(import.meta.url);
const proto = require("./user_pb.cjs");
const PayloadP = proto.toonbench?.PayloadP || proto.PayloadP;
const UserP = proto.toonbench?.UserP || proto.UserP;

const compiled = require("./compiled_protos.cjs");
const toonbench = compiled.toonbench ?? compiled.default ?? compiled;

if (typeof global.gc !== 'function') {
    console.error('Run node with --expose-gc to get reliable memory numbers: node --expose-gc src/bench.js');
    process.exit(1);
}

if (!PayloadP || !UserP) {
  console.error('Não encontrou PayloadP/UserP nas exports de user_pb.js/cjs. Keys:', Object.keys(proto));
  process.exit(1);
}

const NUM_USERS = 5_000;

const roles = ["admin", "user", "moderator", "superuser"]

function generatePayload(n) {
    const users = [];
    for (let i= 0; i < n; i++ ) {
        users.push({
            id: i + 1,
            name: `User ${i+1}`,
            role: roles[i % 4]
        });
    }
    return { users };
}

function generatePayloadProtoGoogle(n) {
    const payload = new PayloadP()
    const users = [];
    for (let i= 0; i < n; i++ ) {
        const user = new UserP();
        user.setId(i + 1);
        user.setName(`User ${i+1}`);
        user.setRole(roles[i % 4]);
        users.push(user)
    }
    payload.setUsersList(users)
    return payload;
}

function generatePayloadProtobufJs(n) {
    const users = [];
    for (let i= 0; i < n; i++ ) {
        const user = toonbench.UserP.create({
            "id": i+1,
            "name": `User ${i+1}`,
            "role": roles[i % 4]
        });
        users.push(user);
    }
    const payload = toonbench.PayloadP.create({
        "users": users
    });
    return payload;
}

const payload  = generatePayload(NUM_USERS);
const jsonStr = JSON.stringify(payload);
const toonStr = encode(payload);

const payloadPGoogle = generatePayloadProtoGoogle(NUM_USERS);
const payloadBufGoogle = payloadPGoogle.serializeBinary();

const payloadPJs = generatePayloadProtobufJs(NUM_USERS);
const payloadBuf = toonbench.PayloadP.encode(payloadPJs).finish();

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

                let nsOpRes = NaN
                if (!Number.isNaN(nsOp)) nsOpRes = Math.round(nsOp).toLocaleString();
                let firstTabs = "\t\t\t";
                if (name.length > 15) {
                    firstTabs = "\t\t";
                }
                if (name.length > 20) {
                    firstTabs = "\t";
                }
                console.log(`${name}${firstTabs}ops/sec: ${hz.toFixed(2)}\tns/op: ${nsOpRes}\texecuted (approx): ${executed}`);
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
        {
            name: 'google-protobuf serializeBin',
            fn: () => payloadPGoogle.serializeBinary()
        },
        {
            name: 'google-protobuf deserializeBin',
            fn: () => PayloadP.deserializeBinary(payloadBufGoogle)
        },
        {
            name: 'protobufjs encode.finish',
            fn: () => toonbench.PayloadP.encode(payloadPJs).finish()
        },
        {
            name: 'protobufjs decode',
            fn: () => toonbench.PayloadP.decode(payloadBuf)
        },
    ];

    const results = [];
    const totalStart = process.hrtime.bigint();

    for (const t of tests) {
        // Run each benchmark isolated to better measure memory impact
        // console.log(`\n--- Running: ${t.name} ---`);
        const r = await runSingleBenchmark(t.name, t.fn);
        results.push(r);
    }

    const totalEnd = process.hrtime.bigint();
    const totalMs = Number(totalEnd - totalStart) / 1e6;

    console.log('\nAll done.');
    console.log(`Total benchmark wall-clock time: ${totalMs.toFixed(0)} ms (${(totalMs/1000).toFixed(3)} s)`);
    // You can further process 'results' to build a table if desired.
}

main().catch(err => {
    console.error("Error:", err);
    process.exit(1);
});
