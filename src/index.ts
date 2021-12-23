import yargs from "yargs";
import { hideBin } from "yargs/helpers";
import { Fs } from "./fs";
import { PathManager } from "./pathManager";
import { registerCommandValidator } from "./validators";
import { getColoredMessage } from "./colors";

const fail = () => process.exit(1);

const main = () => {
  const pathManager = new PathManager(new Fs());
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
    .command(["list", "l"], `list all saved project roots`, {}, async () => {
      await pathManager.listProjects();
    })
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
    .scriptName("")
    .usage("Usage: pr [command]").argv;
};
main();
