import {
  access,
  lstat,
  mkdir,
  readFile,
  truncate,
  writeFile,
} from "fs/promises";
import { EOL, homedir } from "os";
import { join, resolve } from "path";

export class PathManager {
  private readonly storageFile = "storage.txt";
  private storageDir = `${homedir()}/.proot`;

  private get storagePath() {
    return join(this.storageDir, this.storageFile);
  }

  private getSavedPaths(fileContent: string) {
    return fileContent.split(EOL).filter((path) => path.length > 0);
  }

  /**
   *  - Use the database of saved project roots to go to the root directory of the project
   *  - Print the directory to the console such that the bash script can cd to it.
   */
  public async go(cwd: string): Promise<boolean> {
    const fileContent = await this.getFileContent();
    const savedPaths = this.getSavedPaths(fileContent);
    if (savedPaths.length === 0) {
      console.log("There are no registered projects in the databse.");
      return false;
    }
    const pathMatches = savedPaths.filter((path) => cwd.startsWith(path));
    if (pathMatches.length === 0) {
      console.log(
        "The current directory does not belong to a registered project root"
      );
      return false;
    }

    // Take the longest match. should be only 1 path anyway as we are not handling nested roots.
    const gotoPath =
      pathMatches.length === 1
        ? pathMatches[0]
        : pathMatches.sort((a, b) => b.length - a.length)[0];

    console.log(gotoPath);

    return true;
  }

  public async purge(): Promise<boolean> {
    try {
      const fileContent = await this.getFileContent();
      const savedPaths = this.getSavedPaths(fileContent);
      const toDelete: Set<string> = new Set();
      for (const savedPath of savedPaths) {
        const doesExist = await this.doesPathExist(savedPath);
        if (!doesExist) {
          toDelete.add(savedPath);
        }
      }
      if (toDelete.size > 0) {
        const newContent = savedPaths
          .filter((path) => !toDelete.has(path))
          .map((path) => path + EOL)
          .join("");
        await this.writeValue(newContent, false);
        console.log(`Deleted ${toDelete.size} paths:`);
        for (const deletedPath of toDelete) {
          console.log(deletedPath);
        }
      } else {
        console.log("Nothing was deleted");
      }
      return true;
    } catch (error) {
      return false;
    }
  }

  /**
   *
   * register the projectRoot in the list of saved roots.
   */
  public async registerProject(path: string): Promise<boolean> {
    // Make abs path
    if (this.isRelativePath(path)) {
      path = this.toAbsolutePath(path);
    }

    // Check that path exists
    if (!(await this.validatePath(path))) {
      return false;
    }

    // Create storage dir if not exists
    if (!(await this.doesStorageFileExist())) {
      await this.prepareStorageDir();
    }

    const fileContent = await this.getFileContent();

    // Don't handle nested cases.
    const alreadyRegisteredPath = this.tryGetRegisteredPath(fileContent, path);
    if (alreadyRegisteredPath !== undefined) {
      console.log(
        `The path ${path} is already a part of a registered project path ${alreadyRegisteredPath}`
      );
      console.log(
        `To see a list of all registered paths execute the list command`
      );
      return false;
    }

    try {
      await this.writeValue(path + EOL);
    } catch (error) {
      console.log(`Could't write ${path} to ${this.storagePath}`);
      return false;
    }

    console.log(`Successfully registered ${path} as project root.`);

    return true;
  }

  public async listProjects(): Promise<void> {
    const fileContent = await this.getFileContent();
    for (const path of this.getSavedPaths(fileContent)) {
      console.log(path);
    }
  }

  public async clear(): Promise<boolean> {
    try {
      await truncate(this.storagePath, 0);
      console.log("Successfully cleared the database of saved project roots");
      return true;
    } catch {
      console.log(
        "Couldn't clear the database of saved project roots. Does file exist ?"
      );
      return false;
    }
  }

  private async doesStorageFileExist(): Promise<boolean> {
    try {
      await access(this.storagePath);
      return true;
    } catch {
      return false;
    }
  }

  private isRelativePath(path: string) {
    return !path.startsWith("/");
  }

  private toAbsolutePath(path: string) {
    return resolve(path);
  }

  private async doesPathExist(path: string): Promise<boolean> {
    try {
      await access(path);
      return true;
    } catch {
      return false;
    }
  }
  private async validatePath(path: string) {
    try {
      const statResult = await lstat(path);
      if (statResult.isDirectory()) {
        return true;
      } else if (statResult.isFile()) {
        console.log(
          `The project root cannot be a file. It should be a directory`
        );
        return false;
      }
    } catch {
      console.log(`The directory ${path} does not exist in the file system.`);
      return false;
    }
    return false;
  }

  private async getFileContent(): Promise<string> {
    try {
      return (await readFile(this.storagePath)).toString();
    } catch {
      return "";
    }
  }

  /**
   *  Checks if the path to register is withing a registered project.
   */
  private tryGetRegisteredPath(
    storageContent: string,
    pathToRegister: string
  ): string | undefined {
    const savedRoots = this.getSavedPaths(storageContent);
    for (const savedRoot of savedRoots) {
      if (pathToRegister.startsWith(savedRoot)) {
        return savedRoot;
      }
    }
    return undefined;
  }

  private async writeValue(value: string, append = true) {
    await writeFile(this.storagePath, value, {
      encoding: "utf-8",
      flag: append ? "a+" : "w",
    });
  }

  private async prepareStorageDir() {
    await mkdir(this.storageDir, { recursive: true });
  }
}
