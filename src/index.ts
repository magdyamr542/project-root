import yargs from "yargs";
import { hideBin } from "yargs/helpers";
import { Fs } from "./fs";
import { PathManager } from "./pathManager";
import { registerCommandValidator } from "./validators";
import { getColoredMessage } from "./colors";
import { spawn } from "child_process";
import { homedir } from "os";
import { join } from "path";

const fail = () => process.exit(1);

const main = () => {
  const fs = new Fs();
  const pathManager = new PathManager(fs);
  yargs(hideBin(process.argv))
    .command(
      "add [relative path]",
      `register the given path as a root of a project`,
      {},
      async (argv: Record<string, unknown>) => {
        if (registerCommandValidator(argv.relativepath)) {
          const didRegister = await pathManager.registerProject(
            argv.relativepath as string
          );
          if (!didRegister) {
            fail();
          }
        } else {
          fail();
        }
      }
    )
    .command(
      ["list", "l", "ls"],
      `list all saved project roots`,
      {},
      async () => {
        await pathManager.listProjects(process.cwd());
      }
    )
    .command(["go"], `go to the root of this project`, {}, async () => {
      if (!(await pathManager.go(process.cwd()))) {
        fail();
      }
    })
    .command(
      ["clear", "cl"],
      `clear the database of saved projects. Will delete everything. Use with CAUTION`,
      {},
      async () => {
        try {
          if (!(await pathManager.clear())) {
            fail();
          }
        } catch {
          fail();
        }
      }
    )
    .command(
      ["purge", "p"],
      "delete all registered paths that no longer exist in the file system",
      {},
      async () => {
        try {
          if (!(await pathManager.purge())) {
            fail();
          }
        } catch {
          fail();
        }
      }
    )
    .command(
      "$0",
      `the default command. equivalent to pr go.go to the root of this project`,
      () => {},
      async (argv: Record<string, unknown>) => {
        const wrongCommand =
          argv["_"] && Array.isArray(argv["_"]) && argv["_"].length;
        if (wrongCommand) {
          console.log(
            `${getColoredMessage(
              "Error",
              "red",
              true
            )} this command does not exist`
          );
          console.log(
            `${getColoredMessage(
              "Info",
              "grey",
              true
            )} see the command 'pr help'`
          );
          fail();
        }

        if (!(await pathManager.go(process.cwd()))) {
          fail();
        }
      }
    )
    .command(
      ["update"],
      "update the tool (sync with github master branch)",
      {},
      async () => {
        spawn(join(fs.getAppDir(), "update.sh"), { stdio: "inherit" });
      }
    )
    .example(
      "pr",
      `will go to the root of the current project if its path was registered before`
    )
    .example(
      "pr go",
      `will go to the root of the current project if its path was registered before. Same as 'pr'`
    )
    .example(
      "pr add ./",
      `will save the path of the current directory as a project path`
    )
    .example("pr list", "will list all registered paths")
    .example(
      "pr update",
      "will pull the most recent code from github master branch and reinstall the tool"
    )
    .scriptName("")
    .usage("Usage: pr [command]").argv;
};
main();
