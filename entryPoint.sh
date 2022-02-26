function pr {
    output=$(./proot $@)
    retCode=$?
    if [[ ( $@ == "go" || $@ == "" ) && $retCode -eq 0 ]]; then
        # cd when go command. hide output
        cd $output
    else
        echo $output
    fi
    if [ $retCode -ne 0 ]; then
        return $retCode
    fi
    
}
