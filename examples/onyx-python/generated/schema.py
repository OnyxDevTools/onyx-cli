SCHEMA_JSON = {
  "databaseId": "bbabca0e-82ce-11f0-0000-a2ce78b61b6a",
  "entities": [
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "action",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "actorId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "changes",
          "type": "EmbeddedObject"
        },
        {
          "isNullable": True,
          "name": "dateTime",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "errorCode",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "errorMessage",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "metadata",
          "type": "EmbeddedObject"
        },
        {
          "isNullable": True,
          "name": "requestId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "resource",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "status",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "targetId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "tenantId",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "AuditLog",
      "partition": "",
      "resolvers": [],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "deletedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "description",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "name",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "updatedAt",
          "type": "Timestamp"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "Permission",
      "partition": "",
      "resolvers": [],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "deletedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "description",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "isSystem",
          "type": "Boolean"
        },
        {
          "isNullable": True,
          "name": "name",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "updatedAt",
          "type": "Timestamp"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "Role",
      "partition": "",
      "resolvers": [
        {
          "name": "permissions",
          "resolver": "db.from(\"Permission\")\n  .where(\n    inOp(\"id\", \n        db.from(\"RolePermission\")\n            .where(eq(\"roleId\", this.id))\n            .list()\n            .values('permissionId')\n    )\n)\n .list()"
        },
        {
          "name": "rolePermissions",
          "resolver": "db.from(\"RolePermission\")\n .where(eq(\"roleId\", this.id))\n .list()"
        }
      ],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "permissionId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "roleId",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "RolePermission",
      "partition": "",
      "resolvers": [
        {
          "name": "permission",
          "resolver": "db.from(\"Permission\")\n .where(eq(\"id\", this.permissionId))\n .firstOrNull()"
        },
        {
          "name": "role",
          "resolver": "db.from(\"Role\")\n .where(eq(\"id\", this.roleId))\n .firstOrNull()"
        }
      ],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "deletedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "email",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "isActive",
          "type": "Boolean"
        },
        {
          "isNullable": True,
          "name": "lastLoginAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "updatedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "username",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "User",
      "partition": "",
      "resolvers": [
        {
          "name": "profile",
          "resolver": "db.from(\"UserProfile\")\n .where(eq(\"userId\", this.id))\n .firstOrNull()"
        },
        {
          "name": "roles",
          "resolver": "db.from(\"Role\")\n  .where(\n    inOp(\"id\", \n        db.from(\"UserRole\")\n            .where(eq(\"userId\", this.id))\n            .list()\n            .values('roleId')\n    )\n)\n .list()"
        },
        {
          "name": "userRoles",
          "resolver": "db.from(\"UserRole\")\n  .where(eq(\"userId\", this.id))\n  .list()"
        }
      ],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "address",
          "type": "EmbeddedObject"
        },
        {
          "isNullable": True,
          "name": "age",
          "type": "Int"
        },
        {
          "isNullable": True,
          "name": "avatarUrl",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "bio",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "deletedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "firstName",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "lastName",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "phone",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "updatedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "userId",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "UserProfile",
      "partition": "",
      "resolvers": [],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "roleId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "userId",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "UserRole",
      "partition": "",
      "resolvers": [
        {
          "name": "role",
          "resolver": "db.from(\"Role\")\n .where(eq(\"id\", this.roleId))\n .list()"
        }
      ],
      "triggers": [],
      "type": "SEARCHABLE"
    }
  ],
  "revisionDescription": "added full table search",
  "tables": [
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "action",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "actorId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "changes",
          "type": "EmbeddedObject"
        },
        {
          "isNullable": True,
          "name": "dateTime",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "errorCode",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "errorMessage",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "metadata",
          "type": "EmbeddedObject"
        },
        {
          "isNullable": True,
          "name": "requestId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "resource",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "status",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "targetId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "tenantId",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "AuditLog",
      "partition": "",
      "resolvers": [],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "deletedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "description",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "name",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "updatedAt",
          "type": "Timestamp"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "Permission",
      "partition": "",
      "resolvers": [],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "deletedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "description",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "isSystem",
          "type": "Boolean"
        },
        {
          "isNullable": True,
          "name": "name",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "updatedAt",
          "type": "Timestamp"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "Role",
      "partition": "",
      "resolvers": [
        {
          "name": "permissions",
          "resolver": "db.from(\"Permission\")\n  .where(\n    inOp(\"id\", \n        db.from(\"RolePermission\")\n            .where(eq(\"roleId\", this.id))\n            .list()\n            .values('permissionId')\n    )\n)\n .list()"
        },
        {
          "name": "rolePermissions",
          "resolver": "db.from(\"RolePermission\")\n .where(eq(\"roleId\", this.id))\n .list()"
        }
      ],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "permissionId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "roleId",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "RolePermission",
      "partition": "",
      "resolvers": [
        {
          "name": "permission",
          "resolver": "db.from(\"Permission\")\n .where(eq(\"id\", this.permissionId))\n .firstOrNull()"
        },
        {
          "name": "role",
          "resolver": "db.from(\"Role\")\n .where(eq(\"id\", this.roleId))\n .firstOrNull()"
        }
      ],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "deletedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "email",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "isActive",
          "type": "Boolean"
        },
        {
          "isNullable": True,
          "name": "lastLoginAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "updatedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "username",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "User",
      "partition": "",
      "resolvers": [
        {
          "name": "profile",
          "resolver": "db.from(\"UserProfile\")\n .where(eq(\"userId\", this.id))\n .firstOrNull()"
        },
        {
          "name": "roles",
          "resolver": "db.from(\"Role\")\n  .where(\n    inOp(\"id\", \n        db.from(\"UserRole\")\n            .where(eq(\"userId\", this.id))\n            .list()\n            .values('roleId')\n    )\n)\n .list()"
        },
        {
          "name": "userRoles",
          "resolver": "db.from(\"UserRole\")\n  .where(eq(\"userId\", this.id))\n  .list()"
        }
      ],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "address",
          "type": "EmbeddedObject"
        },
        {
          "isNullable": True,
          "name": "age",
          "type": "Int"
        },
        {
          "isNullable": True,
          "name": "avatarUrl",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "bio",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "deletedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "firstName",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "lastName",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "phone",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "updatedAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "userId",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "UserProfile",
      "partition": "",
      "resolvers": [],
      "triggers": [],
      "type": "SEARCHABLE"
    },
    {
      "attributes": [
        {
          "isNullable": True,
          "name": "createdAt",
          "type": "Timestamp"
        },
        {
          "isNullable": True,
          "name": "id",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "roleId",
          "type": "String"
        },
        {
          "isNullable": True,
          "name": "userId",
          "type": "String"
        }
      ],
      "identifier": {
        "generator": "None",
        "name": "id",
        "type": "String"
      },
      "indexes": [],
      "name": "UserRole",
      "partition": "",
      "resolvers": [
        {
          "name": "role",
          "resolver": "db.from(\"Role\")\n .where(eq(\"id\", this.roleId))\n .list()"
        }
      ],
      "triggers": [],
      "type": "SEARCHABLE"
    }
  ]
}
