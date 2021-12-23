import { getColoredMessage } from "./colors";
import { truncate } from "fs/promises";
import { Fs } from "./fs";
import { EOL, homedir } from "os";
export class PathManager {
  private readonly storageFile = "storage.txt";
  private storageDir = `${homedir()}/.proot`;

  constructor(private readonly fs: Fs) {}

  private get storagePath() {
    return this.fs.joinPath(this.storageDir, this.storageFile);
  }

  private getSavedPaths(fileContent: string) {
    return fileContent.split(EOL).filter((path) => path.length > 0);
  }

  get successPrefix() {
    return getColoredMessage("Success", "green");
  }

  get errorPrefix() {
    return getColoredMessage("Error", "red", true);
  }

  get infoPrefix() {
    return getColoredMessage("Info", "grey", true);
  }

  /**
   *  - Use the database of saved project roots to go to the root directory of the project
   *  - Print the directory to the console such that the bash script can cd to it.
   */
  public async go(cwd: string): Promise<boolean> {
    const fileContent = await this.getSavedData();
    const savedPaths = this.getSavedPaths(fileContent);
    if (savedPaths.length === 0) {
      console.log(`${this.errorPrefix} there are no registered projects.`);
      return false;
    }
    const pathMatches = savedPaths.filter((path) => cwd.startsWith(path));
    if (pathMatches.length === 0) {
      console.log(
        `${this.errorPrefix} the current directory does not belong to a registered project`
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
      const fileContent = await this.getSavedData();
      const savedPaths = this.getSavedPaths(fileContent);
      const toDelete: Set<string> = new Set();
      for (const savedPath of savedPaths) {
        const doesExist = await this.fs.doesPathExist(savedPath);
        if (!doesExist) {
          toDelete.add(savedPath);
        }
      }
      if (toDelete.size > 0) {
        const newContent = savedPaths
          .filter((path) => !toDelete.has(path))
          .map((path) => path + EOL)
          .join("");

        await this.fs.writeFile(this.storagePath, newContent, false);
        console.log(`${this.successPrefix} deleted ${toDelete.size} paths:`);
        for (const deletedPath of toDelete) {
          console.log(deletedPath);
        }
      } else {
        console.log(`${this.infoPrefix} nothing was deleted`);
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
    if (this.fs.isRelativePath(path)) {
      path = this.fs.toAbsolutePath(path);
    }

    // Check that path exists
    if (!(await this.validatePath(path))) {
      return false;
    }

    // Create storage dir if not exists
    if (!(await this.fs.doesPathExist(this.storagePath))) {
      await this.prepareStorageDir();
    }

    const fileContent = await this.getSavedData();

    // Don't handle nested cases.
    const alreadyRegisteredPath = this.tryGetRegisteredPath(fileContent, path);
    if (alreadyRegisteredPath !== undefined) {
      console.log(
        `${this.errorPrefix} the path ${path} is already a part of a registered project path ${alreadyRegisteredPath}`
      );
      console.log(
        `${this.infoPrefix} to see a list of all registered paths execute the list command`
      );
      return false;
    }

    try {
      await this.fs.writeFile(this.storagePath, path + EOL);
    } catch (error) {
      console.log(
        `${this.errorPrefix} could't write ${path} to ${this.storagePath}`
      );
      return false;
    }

    console.log(`${this.successPrefix} Added ${path}`);
    return true;
  }

  public async listProjects(): Promise<void> {
    const fileContent = await this.getSavedData();
    for (const path of this.getSavedPaths(fileContent)) {
      console.log(path);
    }
  }

  public async clear(): Promise<boolean> {
    try {
      await truncate(this.storagePath, 0);
      console.log(
        `${this.successPrefix} cleared the database of saved project roots`
      );
      return true;
    } catch {
      console.log(
        `${this.errorPrefix} couldn't clear the database of saved project roots. Does file exist ?`
      );
      return false;
    }
  }

  private async validatePath(path: string) {
    const statResult = await this.fs.statFile(path);
    if (statResult) {
      if (statResult.isDirectory()) {
        return true;
      } else if (statResult.isFile()) {
        console.log(
          `${this.errorPrefix} the project root cannot be a file. It should be a directory`
        );
        return false;
      }
    } else {
      console.log(
        `${this.errorPrefix} the directory ${path} does not exist in the file system.`
      );
      return false;
    }
    return false;
  }

  private async getSavedData(): Promise<string> {
    return this.fs.readFileOrEmptyString(this.storagePath);
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

  private async prepareStorageDir() {
    await this.fs.mkdir(this.storageDir);
  }
}
