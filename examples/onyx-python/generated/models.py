import datetime
from typing import Any, Optional

class AuditLog:
    """Generated model (plain Python class). Resolver/extra fields are allowed via **extra."""
    def __init__(self, id: Optional[str] = None, action: Optional[str] = None, actorId: Optional[str] = None, changes: Optional[dict] = None, dateTime: Optional[datetime.datetime] = None, errorCode: Optional[str] = None, errorMessage: Optional[str] = None, metadata: Optional[dict] = None, requestId: Optional[str] = None, resource: Optional[str] = None, status: Optional[str] = None, targetId: Optional[str] = None, tenantId: Optional[str] = None, **extra: Any):
        self.id = id
        self.action = action
        self.actorId = actorId
        self.changes = changes
        self.dateTime = dateTime
        self.errorCode = errorCode
        self.errorMessage = errorMessage
        self.metadata = metadata
        self.requestId = requestId
        self.resource = resource
        self.status = status
        self.targetId = targetId
        self.tenantId = tenantId
        # allow resolver-attached fields or extra properties
        for k, v in extra.items():
            setattr(self, k, v)


class Permission:
    """Generated model (plain Python class). Resolver/extra fields are allowed via **extra."""
    def __init__(self, id: Optional[str] = None, createdAt: Optional[datetime.datetime] = None, deletedAt: Optional[datetime.datetime] = None, description: Optional[str] = None, name: Optional[str] = None, updatedAt: Optional[datetime.datetime] = None, **extra: Any):
        self.id = id
        self.createdAt = createdAt
        self.deletedAt = deletedAt
        self.description = description
        self.name = name
        self.updatedAt = updatedAt
        # allow resolver-attached fields or extra properties
        for k, v in extra.items():
            setattr(self, k, v)


class Role:
    """Generated model (plain Python class). Resolver/extra fields are allowed via **extra."""
    def __init__(self, id: Optional[str] = None, createdAt: Optional[datetime.datetime] = None, deletedAt: Optional[datetime.datetime] = None, description: Optional[str] = None, isSystem: Optional[bool] = None, name: Optional[str] = None, updatedAt: Optional[datetime.datetime] = None, **extra: Any):
        self.id = id
        self.createdAt = createdAt
        self.deletedAt = deletedAt
        self.description = description
        self.isSystem = isSystem
        self.name = name
        self.updatedAt = updatedAt
        # allow resolver-attached fields or extra properties
        for k, v in extra.items():
            setattr(self, k, v)


class RolePermission:
    """Generated model (plain Python class). Resolver/extra fields are allowed via **extra."""
    def __init__(self, id: Optional[str] = None, createdAt: Optional[datetime.datetime] = None, permissionId: Optional[str] = None, roleId: Optional[str] = None, **extra: Any):
        self.id = id
        self.createdAt = createdAt
        self.permissionId = permissionId
        self.roleId = roleId
        # allow resolver-attached fields or extra properties
        for k, v in extra.items():
            setattr(self, k, v)


class User:
    """Generated model (plain Python class). Resolver/extra fields are allowed via **extra."""
    def __init__(self, id: Optional[str] = None, createdAt: Optional[datetime.datetime] = None, deletedAt: Optional[datetime.datetime] = None, email: Optional[str] = None, isActive: Optional[bool] = None, lastLoginAt: Optional[datetime.datetime] = None, updatedAt: Optional[datetime.datetime] = None, username: Optional[str] = None, **extra: Any):
        self.id = id
        self.createdAt = createdAt
        self.deletedAt = deletedAt
        self.email = email
        self.isActive = isActive
        self.lastLoginAt = lastLoginAt
        self.updatedAt = updatedAt
        self.username = username
        # allow resolver-attached fields or extra properties
        for k, v in extra.items():
            setattr(self, k, v)


class UserProfile:
    """Generated model (plain Python class). Resolver/extra fields are allowed via **extra."""
    def __init__(self, id: Optional[str] = None, address: Optional[dict] = None, age: Optional[int] = None, avatarUrl: Optional[str] = None, bio: Optional[str] = None, createdAt: Optional[datetime.datetime] = None, deletedAt: Optional[datetime.datetime] = None, firstName: Optional[str] = None, lastName: Optional[str] = None, phone: Optional[str] = None, updatedAt: Optional[datetime.datetime] = None, userId: Optional[str] = None, **extra: Any):
        self.id = id
        self.address = address
        self.age = age
        self.avatarUrl = avatarUrl
        self.bio = bio
        self.createdAt = createdAt
        self.deletedAt = deletedAt
        self.firstName = firstName
        self.lastName = lastName
        self.phone = phone
        self.updatedAt = updatedAt
        self.userId = userId
        # allow resolver-attached fields or extra properties
        for k, v in extra.items():
            setattr(self, k, v)


class UserRole:
    """Generated model (plain Python class). Resolver/extra fields are allowed via **extra."""
    def __init__(self, id: Optional[str] = None, createdAt: Optional[datetime.datetime] = None, roleId: Optional[str] = None, userId: Optional[str] = None, **extra: Any):
        self.id = id
        self.createdAt = createdAt
        self.roleId = roleId
        self.userId = userId
        # allow resolver-attached fields or extra properties
        for k, v in extra.items():
            setattr(self, k, v)


