function pr {
    output=$(node /Users/apple/Desktop/frontend_apps/project-root/dist/index.js $@)
    retCode=$?
    if [[ $@ == "go" && $retCode -eq 0 ]]; then
        # cd when go command. hide output
        cd $output
    else
        echo $output
    fi
    if [ $retCode -ne 0 ]; then
       return $retCode
    fi

}
