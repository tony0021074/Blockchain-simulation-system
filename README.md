# COMP5311-GroupProject  
Simluate blockchain  
  
Only database.go is my work. All other files are written by teammates.

The database.go in folder blockchain is the key file in this webpage.  
  
Try the database functionality with all the files uploaded.  
  
Be careful of the .go file directory and the package name  
Modify package name if you move database.go file  
  
2 functions are provided to be used for accessing the blockchain database.  
SaveBlock(block *Block, nodeID string) error ---Create a database and store, or store a block to the existing database  
LoadChain(nodeID string) ([]*Block, error) ---Load all blocks (ordering depends whether the chain is complete) stored in the node's database  
  
  
update 2018-11-27  
1.fix bugs of not loading DB data  
2.reset badger logger  
3.sort blocks with bubble sort  
