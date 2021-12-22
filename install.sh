mkdir -p ~/.proot
npm i
npm run build
npm run mkExec
cp ./bin/index.js ./entryPoint.sh  ~/.proot
echo Installed. Please add this line to your .bashrc or .zshrc  "'source .proot/entryPoint.sh'"
