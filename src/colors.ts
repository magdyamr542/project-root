import { bold, green, grey, red } from "colors";

export type MessageColor = "green" | "grey" | "red";
const makeBold = (text: string) => bold(text);
export const getColoredMessage = (
  text: string,
  color: MessageColor,
  bold = false
): string => {
  let message = "";
  if (color === "red") {
    message = red(text);
  } else if (color === "grey") {
    message = grey(text);
  } else {
    message = green(text);
  }
  return bold ? makeBold(message) : message;
};
