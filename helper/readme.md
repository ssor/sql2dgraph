
# Instruction

This app is used to show how to use dgraph when mutate data and schema.

If data need mutation, make sure if new or just update.

If data should be updated to exists node, there should be an uid, or new node will be created, and as a result, data is not in its position.

So there is two method in helper, MutationObj and UpdateObj which are almost same, but latter will check for uid to make sure no new node is created.



