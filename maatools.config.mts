import type {FullConfig} from "@nekosu/maa-tools";

const config: FullConfig = {
    cwd: import.meta.dirname,
    maaVersion: "v4.0.0-beta.14",
    interfacePath: "assets/interface.json",
    check: {},
    vscode: {
        agents: {
            "agent/go-service": "launch-go-agent",
        },
    },
};

export default config;
