## shell script to compile sql files in a folder into a single file

folder=$1
output=$2

if [ -z "$folder" ]; then
    echo "Please provide a folder name"
    exit 1
fi

if [ ! -d "$folder" ]; then
    echo "Folder $folder does not exist"
    exit 1
fi

if [ -z "$output" ]; then
    echo "Please provide an output file name"
    exit 1
fi

for file in $folder/*.pgsql; do
    echo "Compiling $file"
    cat $file >> $output
done