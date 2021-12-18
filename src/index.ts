import { EOL } from "os";
import yargs from "yargs";
import { hideBin } from "yargs/helpers";
import { PathManager } from "./pathManager";
import { registerCommandValidator } from "./validators";

const fail = () => process.exit(1);

const main = () => {
  const pathManager = new PathManager();
  yargs(hideBin(process.argv))
    .command(
      "register [relative path]",
      `register the given path as a root of a project${EOL}`,
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
    .command("list", `list all saved project roots${EOL}`, {}, async () => {
      await pathManager.listProjects();
    })
    .command("go", `go to the root of this project${EOL}`, {}, async () => {
      if (!(await pathManager.go(process.cwd()))) {
        fail();
      }
    })
    .command(
      "clear",
      `clear the database of saved projects. Will delete everything. Use with CAUTION${EOL}`,
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
      "purge",
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
    .example(
      "proot ./",
      "will save the path of the current directory as a project path"
    )
    .example(
      "proot go",
      "will go to the root of the current project if its path was registered before"
    )
    .scriptName("proot")
    .usage("Usage: proot [command]").argv;
};
main();
