echo "Updating project-root..."
mkdir -p ~/.proot/updating_project
cd ~/.proot/updating_project
git clone git@github.com:magdyamr542/project-root.git
cd project-root 
./install.sh
rm -rf ~/.proot/updating_project
echo "Updated"


