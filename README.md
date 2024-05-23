# Go ECS

A very minimal Entity Component System implementation I created to help me learn more about Go (especially the new generics) and ECS architectures.

# Concepts

## Buckets
Buckets are collections of entities and their associated components

A bucket contains a map of component types (string representations) to slices of component "entries". The length of these slices is equal to the number of entities that have been added to the bucket (entities can be "deleted" from a bucket but their index into component slices is just reused, the component slices lengths are never reduced, but the can grow).


## Entities
Entities in this ECS are just containers that hold indexes into the component slices in the bucket to which they are added and a bit mask indicating which component types they have associated with them.


## Components
Components can take any form and are just pointers to `interface{}` values and are stored in their respective component slice in the bucket to which they are added.

## Queries
Queries are basically just a container that holds a bit mask that is updated when new component types are added to the query. The querie's bit mask is compared with all the entities in the bucket when the query is executed and only the entities whose own masks match the query mask are included in the results.

The components from a query result are accessible in a type safe manner via type assertions and these assertions are abstracted away behind the easy to use `GetComponentFromQueryResult` function.

## Systems
Systems aren't really defined in this simple ECS implementation, but queries are, so systems can be functions that execute a query and then act on the query results or anything else.


