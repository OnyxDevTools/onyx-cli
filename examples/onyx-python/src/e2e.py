
from onyx_database import onyx
from onyx import SCHEMA, tables, User


def main() -> None:
    db = onyx.init(schema=SCHEMA)

    # CRUD via Python SDK using generated onyx
    record_id = "cli-py-e2e"
    new_user = User(id=record_id, name="Py E2E")

    print("Creating User...")
    created = db.save(tables.User, new_user)
    print("Create response:", created)

    fetched = db.find_by_id(tables.User, record_id)
    print("Get response:", fetched)

    updated = db.save(tables.User, {"id": record_id, "name": "Py E2E Updated"})
    print("Update response:", updated)

    deleted = db.delete(tables.User, record_id)
    print("Delete response:", deleted)

    print("Python example CLI+SDK compatibility test passed.")


if __name__ == "__main__":
    main()
