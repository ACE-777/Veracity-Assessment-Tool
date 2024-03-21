# Configuration
## For local execution
### Windows
Set env variables in system settings for path with 
usage python interpreter for python repo with colored
and usage conda. Then reboot computer.
![img.png](img.png)
## Versions of coloring
branches:
- dev/colored/v1/1 - coloring first version
- dev/colored/v2/1 - coloring second version with new approache without skipping

soon...:
- dev/colored/v1/2 - coloring with lucine
- dev/colored/v2/2 - coloring with skiping 3 tokens 
- dev/colored/v2/3 - coloring with colBert  

## To generate proto file
protoc -I=. --go_out=. ./matrix.proto
