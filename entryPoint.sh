function pr {
    output=$(~/.proot/proot $@)
    retCode=$?
    if [[ ( $1 == "go" || $1 == "" || $1 == "back" || $1 == "b" || $1 == "to" || $1 == "t" ) && $retCode -eq 0 ]]; then
        # cd when go or back command. hide output
        cd $output
    else
        echo $output
    fi
    if [ $retCode -ne 0 ]; then
        return $retCode
    fi
    
}
