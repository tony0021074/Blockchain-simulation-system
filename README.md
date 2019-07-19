# COMP5311-GroupProject  
The application aims to simluate the mechanism of blockchain. It is a group project. I am responsible to the database funtionality. The database.go in folder blockchain is the file written by me.  
  
Try the database functionality with all the files uploaded.  
  
Be careful of the .go file directory and the package name  
Modify package name if you move database.go file  

3 functions are provided to be used for accessing the blockchain database:  
SaveNode(nodeID string) ---Create a database for the node  
SaveBlockToNode(nodeID string, block *Block) ---Store a block to the node's database  
LoadBlocksFromNode(nodeID string) ---Load all blocks stored in the node's database  
