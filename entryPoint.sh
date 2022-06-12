function pr {
    output=$(~/.proot/proot $@)
    retCode=$?
    if [[ ( $@ == "go" || $@ == "" || $@ == "back" || $@ == "b" || $@ == "to" || $@ == "t" ) && $retCode -eq 0 ]]; then
        # cd when go or back command. hide output
        cd $output
    else
        echo $output
    fi
    if [ $retCode -ne 0 ]; then
        return $retCode
    fi
    
}
