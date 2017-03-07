## linked-data

This module defines a linked-data validation for data models in envoke. When a model contains the ids of other models, those models are queried and validated to determine its validity. For instance, a publication contains a composition id, and a composition contains a composer id. To validate the publication, we must check that the composer of the linked composition is a valid agent in the database. We can think of this as a graph and the validation process a traversal of the graph; we follow edges to nodes (data models) and check whether they have certain properties. The following sections summarize the core validation processes. For details on data models and field validation, check out the `spec`.

![linked_data](https://github.com/zbo14/envoke/blob/master/linked_data/images/linked_data.png?raw=true)

### Composition

- Validate fields in composition
- Query composer, publisher and validate fields
- Check that composer signed composition tx

### Composition Right Assignment 

- Validate fields in assignment 
- Query holder, signer and validate fields
- Check that signer signed assignment tx 
- Query composition right and do the following...
	- Check that holder holds primary output of right tx 
	- Check that signer signed right tx

### Publication

- Validate fields in publishing license
- Query composition and run validation process
- Check that publisher signed tx
- For each composition right assignment...
	- Query and run validation process
	- Check that underlying right links to composition
	- Check that previous right assignments do not link to underlying right
	- Check that composer is signer of right assignment
	- Check that holder does not hold another right assignment to composition
- Check that total percentage shares from underlying rights equal 100

### Mechanical License

- Validate fields in mechanical license
- Query licensee, licenser and validate fields
- Check that licenser signed license tx
- Query publication and run validation process
- Check that publication links to composition right assignment linked to in license
- Check that licenser is holder of right assignment
- Check that license territory is subset of right territory, if territory specified

### Composition Right Transfer

- Validate fields in transfer
- Query recipient, sender and validate fields
- Check that recipient and sender have different keys 
- Query publication and run validation process
- Query TRANSFER tx and do the following...
	- Check that it has TRANSFER operation
	- Check that it was signed by sender
	- Check that recipient holds primary output and sender holds secondary output, if there is one
- Check that publication links to right assignment that links to composition right linked to in TRANSFER tx

### Recording

- Validate fields in the recording
- Query label, performer, producer and validate fields
- Check that performer/producer signed recording tx
- Query publication and run validation process
- If recording links to a composition right assignment...
	- Check that publication links to it
	- Check that recording signer is holder of right assignment

### Recording Right Assignment 

- Validate fields in assignment 
- Query holder, signer and validate fields
- Check that signer signed tx 
- Query recording right and do the following...
	- Check that holder holds primary output of tx 
	- Check that signer signed tx

### Release

- Validate fields in the release
- Query recording and run the validation process
- Check that recording label signed release tx
- If recording does not link to composition right assignment...
	- Query mechanical license and run validation process
	- Check that mechanical license links to correct publication
	- Check that recording label is licensee
- For each recording right assignment...
	- Query and run validation process
	- Check that underlying right links to recording
	- Check that previous right assignments do not link to underlying right
	- Check that recording signer is signer of right assignment
	- Check that holder does not hold another right assignment to recording
- Check that total percentage shares from underlying rights equal 100 

### Master License

- Validate fields in master license
- Query licensee, licenser and validate fields
- Check that licenser signed license tx
- Query release and run validation process
- Check that release links to recording right assignment linked to in license
- Check that licenser is holder of right assignment
- Check that license territory is subset of right territory, if territory specified

### Recording Right Transfer

- Validate fields in transfer
- Query recipient, sender and validate fields
- Check that recipient and sender have different keys 
- Query release and run validation process
- Query TRANSFER tx and do the following...
	- Check that it has TRANSFER operation
	- Check that it was signed by sender
	- Check that recipient holds primary output and sender holds the secondary output, if there is one
- Check that release links to right assignment that links to recording right linked to in TRANSFER tx