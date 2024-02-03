#!/usr/bin/env -S deno run -A
import { Tar } from "https://deno.land/std@0.184.0/archive/mod.ts";
import { copy } from "https://deno.land/std@0.184.0/streams/copy.ts";

const GOARCH = [
    "amd64",
    "arm64",
];

const GOOS = [
    "windows",
    "linux",
    "darwin",
    "freebsd",
    "netbsd",
    "openbsd",
    "plan9",
    "solaris",
];

await Deno.mkdir("./build", { recursive: true });

console.time("Built ovm");
Deno.env.set("CGO_ENABLED", "0");

for (const os of GOOS) {
    for (const ar of GOARCH) {
        if (os == "solaris" && ar == "arm64" || os == "plan9" && ar == "arm64") {
            continue;
        }
        Deno.env.set("GOOS", os);
        Deno.env.set("GOARCH", ar);
        const ovm_str = `ovm-${os}-${ar}`;
        console.time(`Build ovm: ${ovm_str}`);
        // deno-lint-ignore no-deprecated-deno-api
        const build_cmd = Deno.run({
            cmd: [
                "go",
                "build",
                "-o",
                `build/${ovm_str}/ovm${(os == "windows" ? ".exe" : "")}`,
                "-ldflags=-w -s", "-trimpath",
            ],
        });

        const { code } = await build_cmd.status();
        if (code !== 0) {
            console.error("Something went wrong");
            Deno.exit(1);
        }

        console.timeEnd(`Build ovm: ${ovm_str}`);
    }
}

Deno.chdir("build");
for (const os of GOOS) {
    for (const ar of GOARCH) {
        if (os == "solaris" && ar == "arm64" || os == "plan9" && ar == "arm64") {
            continue;
        }
        const ovm_str = `ovm-${os}-${ar}`;

        if (os == "windows") {
            console.time(`Compress ovm: ${ovm_str}`);
            const zip = new Deno.Command(`zip`, {
                args: [`${ovm_str}.zip`, `${ovm_str}/ovm.exe`],
                stdin: "piped",
                stdout: "piped",
            });
            zip.spawn();
            console.timeEnd(`Compress ovm: ${ovm_str}`);
            continue;
        }
        const tar = new Tar();
        console.time(`Compress ovm: ${ovm_str}`);
        await tar.append("ovm", {
            filePath: `${ovm_str}/ovm`,
        });
        const writer = await Deno.open(`./${ovm_str}.tar`, {
            write: true,
            create: true,
        });
        await copy(tar.getReader(), writer);
        writer.close();
        console.timeEnd(`Compress ovm: ${ovm_str}`);
    }
}

console.timeEnd(`Built ovm`);
