## linked-data

This module defines a linked-data validation for data models in envoke. When a model contains the ids of other models, those models must be queried and validated to determine its validity. For instance, a recording contains a composition id and a composition contains a composer id. To validate the recording, we must check that the composer of the composition which this is a recording of, is a valid agent in the database. We can think of this as a graph and the validation process a traversal of the graph; we follow edges to nodes (data models) and check whether they have certain properties. The following sections summarize the core validation processes. For details on data models and field validation, check out the `spec`.

### Publishing

![publishing](https://github.com/zbo14/envoke/blob/master/linked_data/images/publishing.png?raw=true)

(1) Validation process for a composition

- Validate fields in composition
- Query composer, publisher and validate fields
- Check that composition key matches composer or publisher key
- Query each right holder and validate fields

(2) Validation process for a publishing license

- Check that publishing license has valid fields
- Query composition and run validation process (1)
- Check that licenser id matches a composition right holder id
- Query licenser and check that key matches license key (*already validated fields)
- Query licensee and validate fields

### Recording

![recording](https://github.com/zbo14/envoke/blob/master/linked_data/images/recording.png?raw=true)

(3) Validation process for a recording 

- Validate fields in the recording
- Query composition and run validation process (1)
- Query label, perforner, producer and validate fields
- Check that recording key matches label or performer key
- Agent with matching key must hold composition right or publishing license
- If there's a publishing license, run the validation process (2)
- Query each right holder and validate fields

(4) Validation process for a recording license

- Validate fields in the recording license
- Query the recording and run the validation process (3)
- Check that the licenser id matches a recording right holder id
- Query the licenser and check that key matches license key
- Query the licensee and validate fields 

