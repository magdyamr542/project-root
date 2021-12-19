import { Stats } from "fs";
import { access, lstat, mkdir, readFile, writeFile } from "fs/promises";
import { join, resolve } from "path";

export class Fs {
  public joinPath(...paths: string[]) {
    return join(...paths);
  }

  public async statFile(filePath: string): Promise<Stats | undefined> {
    try {
      return await lstat(filePath);
    } catch (error) {
      return undefined;
    }
  }

  public async doesPathExist(filePath: string): Promise<boolean> {
    try {
      await access(filePath);
      return true;
    } catch {
      return false;
    }
  }

  public isRelativePath(path: string) {
    return !path.startsWith("/");
  }

  public toAbsolutePath(path: string) {
    return resolve(path);
  }

  public async readFileOrEmptyString(filePath: string): Promise<string> {
    try {
      return (await readFile(filePath)).toString();
    } catch {
      return "";
    }
  }

  public async writeFile(filePath: string, value: string, append = true) {
    await writeFile(filePath, value, {
      encoding: "utf-8",
      flag: append ? "a+" : "w",
    });
  }

  public async mkdir(dirPath: string) {
    await mkdir(dirPath, { recursive: true });
  }
}
