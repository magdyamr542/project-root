export const registerCommandValidator = (
  path: string | unknown | undefined
): boolean => {
  if (path === undefined) {
    console.log("Please enter a relative path to the project root");
    return false;
  } else if (typeof path !== "string") {
    console.log("Please enter a relative path as a string");
    return false;
  } else {
    return true;
  }
};
