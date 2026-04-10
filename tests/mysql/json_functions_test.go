package mysql

import (
	"encoding/json"
	"testing"

	"github.com/go-jet/jet/v2/internal/testutils"
	. "github.com/go-jet/jet/v2/mysql"
	. "github.com/go-jet/jet/v2/tests/.gentestdata/mysql/dvds/table"
	"github.com/stretchr/testify/require"
)

func TestJSON_OBJECT_Function(t *testing.T) {
	onlyMariaDB(t)

	stmt := SELECT(
		JSON_OBJECT(
			Actor.ActorID.AS("actorID"),
			Actor.FirstName.AS("firstName"),
		).AS("json"),
	).FROM(Actor).
		WHERE(Actor.ActorID.EQ(Int(2)))

	testutils.AssertStatementSql(t, stmt, `
SELECT JSON_OBJECT(
          'actorID', actor.actor_id,
          'firstName', actor.first_name
     ) AS "json"
FROM dvds.actor
WHERE actor.actor_id = ?;
`, int64(2))

	var dest struct {
		Json string
	}

	err := stmt.Query(db, &dest)
	require.NoError(t, err)
	require.JSONEq(t, `{"actorID": 2, "firstName": "NICK"}`, dest.Json)
	requireLogged(t, stmt)
	requireQueryLogged(t, stmt, 1)
}

func TestJSON_ARRAYAGG_Function(t *testing.T) {
	onlyMariaDB(t)

	stmt := SELECT(
		JSON_ARRAYAGG(
			Actor.ActorID.AS("actorID"),
		).ORDER_BY(Actor.ActorID.ASC()).
			LIMIT(3).
			AS("json"),
	).FROM(Actor)

	testutils.AssertStatementSql(t, stmt, `
SELECT JSON_ARRAYAGG(JSON_OBJECT(
          'actorID', actor.actor_id
     )
     ORDER BY actor.actor_id ASC
     LIMIT ?) AS "json"
FROM dvds.actor;
`, int64(3))

	var dest struct {
		Json string
	}

	err := stmt.Query(db, &dest)
	require.NoError(t, err)
	require.JSONEq(t, `[{"actorID":1},{"actorID":2},{"actorID":3}]`, dest.Json)
	requireLogged(t, stmt)
	requireQueryLogged(t, stmt, 1)
}

func TestJSON_ARRAYAGG_FunctionDistinct(t *testing.T) {
	onlyMariaDB(t)

	stmt := SELECT(
		JSON_ARRAYAGG(
			Actor.LastName.AS("lastName"),
		).DISTINCT().
			ORDER_BY(Actor.LastName.ASC()).
			AS("json"),
	).FROM(Actor)

	testutils.AssertStatementSql(t, stmt, `
SELECT JSON_ARRAYAGG(DISTINCT JSON_OBJECT(
          'lastName', actor.last_name
     )
     ORDER BY actor.last_name ASC) AS "json"
FROM dvds.actor;
`)

	var dest struct {
		Json string
	}

	err := stmt.Query(db, &dest)
	require.NoError(t, err)

	var expected []struct {
		LastName string
	}

	err = SELECT(
		Actor.LastName.AS("lastName"),
	).FROM(Actor).
		GROUP_BY(Actor.LastName).
		ORDER_BY(Actor.LastName.ASC()).
		Query(db, &expected)
	require.NoError(t, err)

	expectedJSON := make([]map[string]string, len(expected))
	for i, row := range expected {
		expectedJSON[i] = map[string]string{"lastName": row.LastName}
	}

	expectedJSONBytes, err := json.Marshal(expectedJSON)
	require.NoError(t, err)

	require.JSONEq(t, string(expectedJSONBytes), dest.Json)
	requireLogged(t, stmt)
	requireQueryLogged(t, stmt, 1)
}

func TestJSON_ARRAYAGG_FunctionLimitOffset(t *testing.T) {
	onlyMariaDB(t)

	stmt := SELECT(
		JSON_ARRAYAGG(
			Actor.ActorID.AS("actorID"),
		).ORDER_BY(Actor.ActorID.ASC()).
			LIMIT(3).
			OFFSET(2).
			AS("json"),
	).FROM(Actor)

	testutils.AssertStatementSql(t, stmt, `
SELECT JSON_ARRAYAGG(JSON_OBJECT(
          'actorID', actor.actor_id
     )
     ORDER BY actor.actor_id ASC
     LIMIT ?
     OFFSET ?) AS "json"
FROM dvds.actor;
`, int64(3), int64(2))

	var dest struct {
		Json string
	}

	err := stmt.Query(db, &dest)
	require.NoError(t, err)
	require.JSONEq(t, `[{"actorID":3},{"actorID":4},{"actorID":5}]`, dest.Json)
	requireLogged(t, stmt)
	requireQueryLogged(t, stmt, 1)
}

func TestJSON_OBJECT_FunctionNullValue(t *testing.T) {
	onlyMariaDB(t)

	stmt := SELECT(
		JSON_OBJECT(
			Film.FilmID.AS("filmID"),
			Film.OriginalLanguageID.AS("originalLanguageID"),
		).AS("json"),
	).FROM(Film).
		WHERE(Film.OriginalLanguageID.IS_NULL()).
		LIMIT(1)

	testutils.AssertStatementSql(t, stmt, `
SELECT JSON_OBJECT(
          'filmID', film.film_id,
          'originalLanguageID', film.original_language_id
     ) AS "json"
FROM dvds.film
WHERE film.original_language_id IS NULL
LIMIT ?;
`, int64(1))

	var dest struct {
		Json string
	}

	err := stmt.Query(db, &dest)
	require.NoError(t, err)

	var got map[string]interface{}
	err = json.Unmarshal([]byte(dest.Json), &got)
	require.NoError(t, err)

	_, exists := got["originalLanguageID"]
	require.True(t, exists)
	require.Nil(t, got["originalLanguageID"])
	requireLogged(t, stmt)
	requireQueryLogged(t, stmt, 1)
}

func TestJSON_ARRAYAGG_FunctionGroupBy(t *testing.T) {
	onlyMariaDB(t)

	stmt := SELECT(
		Customer.StoreID.AS("storeID"),
		JSON_ARRAYAGG(
			Customer.CustomerID.AS("customerID"),
		).ORDER_BY(Customer.CustomerID.ASC()).AS("customers"),
	).FROM(Customer).
		GROUP_BY(Customer.StoreID).
		ORDER_BY(Customer.StoreID.ASC())

	testutils.AssertStatementSql(t, stmt, `
SELECT customer.store_id AS "storeID",
     JSON_ARRAYAGG(JSON_OBJECT(
          'customerID', customer.customer_id
     )
     ORDER BY customer.customer_id ASC) AS "customers"
FROM dvds.customer
GROUP BY customer.store_id
ORDER BY customer.store_id ASC;
`)

	var dest []struct {
		StoreID   int16
		Customers string
	}

	err := stmt.Query(db, &dest)
	require.NoError(t, err)

	var expectedCounts []struct {
		StoreID int16
		Cnt     int64
	}

	err = SELECT(
		Customer.StoreID.AS("storeID"),
		COUNT(Customer.CustomerID).AS("cnt"),
	).FROM(Customer).
		GROUP_BY(Customer.StoreID).
		ORDER_BY(Customer.StoreID.ASC()).
		Query(db, &expectedCounts)
	require.NoError(t, err)
	require.Len(t, dest, len(expectedCounts))

	for i, row := range dest {
		require.Equal(t, expectedCounts[i].StoreID, row.StoreID)

		var customers []struct {
			CustomerID int16 `json:"customerID"`
		}

		err = json.Unmarshal([]byte(row.Customers), &customers)
		require.NoError(t, err)
		require.Len(t, customers, int(expectedCounts[i].Cnt))

		for j := 1; j < len(customers); j++ {
			require.LessOrEqual(t, customers[j-1].CustomerID, customers[j].CustomerID)
		}
	}

	requireLogged(t, stmt)
	requireQueryLogged(t, stmt, int64(len(dest)))
}
