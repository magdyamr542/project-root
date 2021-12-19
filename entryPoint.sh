output=$(node /Users/apple/Desktop/frontend_apps/project-root/dist/index.js $@)
retCode=$?
# cd when go command
if [[ $@ == "go" && $retCode -eq 0 ]]; then
    
    cd $output
fi
echo $output
if [ $retCode -ne 0 ]; then
    exit 1
fi
