## linked-data

This module defines a linked-data validation for data models in envoke. When a model contains the ids of other models, those models are queried and validated to determine its validity. For instance, a publication contains a composition id, and a composition contains a composer id. To validate the publication, we must check that the composer of the linked composition is a valid agent in the database. We can think of this as a graph and the validation process a traversal of the graph; we follow edges to nodes (data models) and check whether they have certain properties. The following sections summarize the core validation processes. For details on data models and field validation, check out the `spec`.

[linked_data](https://github.com/zbo14/envoke/blob/master/linked_data/images/linked_data.png?raw=true)

### Composition

- Validate fields in composition
- Query composer, publisher and validate fields
- Check that composer signed tx

### Publication

- Validate fields in publishing license
- Query composition and run validation process
- Check that publisher signed tx
- For each composition right...
	- Query and valdiate fields 
	- Check that right links to composition
	- Check that composer signed right
	- Check that right-holder does not hold another right to composition
- Check that total percentage shares equal 100

### Mechanical License

- Validate fields in mechanical license
- Query licensee, licenser and validate fields
- Check that licenser signed tx
- Query publication and run validation process
- Check that publication links to license composition right
- Query right and check that licenser is right-holder
- Check that license territory is subset of right territory

### Recording

- Validate fields in the recording
- Query label, performer, producer and validate fields
- Check that performer/producer signed tx
- Query publication and run validation process
- If recording links to a composition right...
	- Check that publication links to it
	- Check that recording signer is right-holder

### Release

- Validate fields in the release
- Query recording and run the validation process
- Check that recording label signed tx
- If recording does not link to composition right...
	- Query mechanical license and run validation process
	- Check that mechanical license links to correct publication
	- Check that recording label is licensee
- For each recording right...
	- Query and validate fields
	- Check that right links to recording
	- Check that recording signer signed right
	- Check that right-holder does not hold another right to recording 
- Check that total percentage shares equal 100 

### Master License

- Validate fields in master license
- Query licensee, licenser and validate fields
- Check that licenser signed tx
- Query release and run validation process
- Check that release links to license recording right
- Query right and check that licenser is right-holder
- Check that license territory is subset of right territory