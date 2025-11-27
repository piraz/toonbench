import time
import gc
import json
import tracemalloc
import argparse
from typing import Callable, List, Dict, Any
from toon_format import encode as toon_encode, decode as toon_decode
import orjson
import ujson
from toonbench.user_pb2 import PayloadP, UserP


NUM_USERS = 5_000
DEFAULT_TIME = 1.0
ROLES = ["admin", "user", "moderator", "superuser"]


def generate_payload(n: int) -> Dict[str, Any]:
    users = []
    for i in range(n):
        users.append({
            "id": i + 1,
            "name": f"User_{i}",
            "role": ROLES[i % len(ROLES)]
        })
    return users


def generate_payload_protobuf(n: int) -> PayloadP:
    payload = PayloadP()
    for i in range(n):
        user = UserP()
        user.id = i + 1
        user.name = f"User_{i}"
        user.role = ROLES[i % len(ROLES)]
        payload.users.append(user)
    return payload


def bench_loop(fn: Callable[[], None],
               min_seconds: float = 1.0) -> (int, float):
    start = time.perf_counter()
    executed = 0
    while True:
        fn()
        executed += 1
        now = time.perf_counter()
        if (now - start) >= min_seconds:
            return executed, now - start


def format_mb(bytes_num: int) -> str:
    return f"{bytes_num / 1024 / 1024:.3f} MB"


def run_one(name: str, fn: Callable[[], None],
            min_seconds: float = 1) -> Dict[str, Any]:
    gc.collect()
    tracemalloc.start()
    before_curr, before_peak = tracemalloc.get_traced_memory()

    executed, elapsed = bench_loop(fn, min_seconds=min_seconds)

    after_curr, after_peak = tracemalloc.get_traced_memory()
    tracemalloc.stop()
    gc.collect()

    mem_delta_curr = max(0, after_curr - before_curr)
    mem_delta_peak = max(0, after_peak - before_peak)

    hz = executed / elapsed if elapsed > 0 else float("nan")
    ns_per_op = (elapsed / executed) * 1e9 if executed > 0 else float("nan")
    # elapsed_ms = elapsed * 1000.0

    # print(
    #         f"{name:28} ops/sec: {hz:10.2f}\tns/op: {int(ns_per_op):10,}\t"
    #         f"executed(approx): {executed}\telapesed(ms): {elapsed_ms:.0f}"
    # )
    #
    # print(
    #         f"  mem delta (current): {mem_delta_curr} bytes "
    #         f"({format_mb(mem_delta_curr)}), peak-delta: {mem_delta_peak} "
    #         f"bytes ({format_mb(mem_delta_peak)})"
    # )

    return {
        "name": name,
        "hz": hz,
        "ns_op": ns_per_op,
        "executed": executed,
        "elapsed": elapsed,
        "mem_delta_curr": mem_delta_curr,
        "mem_delta_peak": mem_delta_peak,
    }


def main():
    parser = argparse.ArgumentParser(
            description="TOON-focused Python microbenchmark")
    parser.add_argument("--users", type=int,
                        default=NUM_USERS,
                        help="number of users in payload")
    parser.add_argument("--time", type=float, default=DEFAULT_TIME,
                        help="min seconds per test")
    parser.add_argument("--skip-json", action="store_true",
                        help="skip JSON tests")
    parser.add_argument("--skip-toon", action="store_true",
                        help="skip TOON tests")
    parser.add_argument("--skip-proto", action="store_true",
                        help="skip Protobuf tests")
    args = parser.parse_args()

    print("TOON Python microbenchmark")
    print("NUM_USERS:", args.users)
    print("min seconds per test:", args.time)
    print()

    payload = generate_payload(args.users)
    json_str = json.dumps(payload)
    payload_proto = generate_payload_protobuf(args.users)
    proto_buf = payload_proto.SerializeToString()

    toon_buf = toon_encode(payload)

    tests: List[tuple] = []

    if not args.skip_json:
        tests.append(("JSON.dumps", lambda: json.dumps(payload)))
        tests.append(("JSON.loads", lambda: json.loads(json_str)))
        tests.append(("orjson.dumps", lambda: orjson.dumps(payload)))
        tests.append(("orjson.loads", lambda: orjson.loads(json_str)))
        tests.append(("ujson.dumps", lambda: ujson.dumps(payload)))
        tests.append(("ujson.loads", lambda: ujson.loads(json_str)))

    if not args.skip_proto:
        tests.append(("protobuf.SerializeToString",
                      lambda: payload_proto.SerializeToString()))
        tests.append(("protobuf.ParseFromString",
                      lambda: PayloadP().ParseFromString(proto_buf)))

    if not args.skip_toon:
        tests.append(("TOON.encode", lambda: toon_encode(payload)))
        if toon_buf is not None:
            tests.append(("TOON.decode", lambda: toon_decode(toon_buf)))
        else:
            print("Skipping TOON.decode because initial encoding failed or "
                  "was not prepared.")

    if not tests:
        print("No tests to run. Use --skip-json/--skip-toon to adjust. "
              "Exiting.")
        return

    results = []
    total_start = time.perf_counter()
    for name, fn in tests:
        # print(f"\n--- Running: {name} ---")
        r = run_one(name, fn, min_seconds=args.time)
        results.append(r)
    total_elapsed = (time.perf_counter() - total_start) * 1000.0
    print(f"\nAll done. Total wall-clock for Python bench run: "
          f"{total_elapsed:.0f} ms ({total_elapsed/1000.0:.3f} s)\n")

    print("Summary:")
    for r in results:
        print(f"- {r['name']:28} elapsed(ms): {r['elapsed']*1000:.0f}  "
              f"ops/sec: {r['hz']:.2f}  ns/op: {int(r['ns_op']):,}")


if __name__ == "__main__":
    main()
