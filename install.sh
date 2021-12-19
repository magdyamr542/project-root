mkdir -p ~/.proot
npm i
npm run build
npm run mkbin
cp ./bin/index.js ./entryPoint.sh ~/.proot
if [[ -f  ~/.bashrc ]]; then
    echo "source ~/.proot/entryPoint.sh"  >> ~/.bashrc
fi

if [[ -f  ~/.zshrc ]]; then
    echo "source ~/.proot/entryPoint.sh"  >> ~/.zshrc
fi

echo "The app was installed succesffully. Source your shell again and use pr help to get started"