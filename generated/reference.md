# Reference
## Users Management Version 2.0
<details><summary><code>client.UsersManagementVersion20.GetUsers() -> *generated.GetUsersResponseV2</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlement required: readUserGroups (Read users and groups), or manageUserGroups (Manage users and groups), or manageAllUserGroups (Synchronize users and groups), or manageUserStandardGroups (Manage users and standard groups).<br>Note: You only need one entitlement, but you can have more than one.<br><br>Searching, sorting, paging, and filtering are supported. A maximum of 2500 records are returned for a search query.<br><br>To improve performance, specify the list of the attributes that you want returned by using the attributes query parameter.<br><br>Search operators supported:<table><tr><td>eq</td><td>The attribute and operator values must be identical for a match.</td></tr><tr><td>ne</td><td>The attribute and operator values are not identical.</td></tr><tr><td>co</td><td>The entire operator value must be a substring of the attribute value for a match.  For performance reasons, use sw or ew operators instead of co.</td></tr><tr><td>sw</td><td>The entire operator value must be a substring of the attribute value, starting at the beginning of the attribute value.</td></tr><tr><td>ew</td><td>The entire operator value must be a substring of the attribute value, matching at the end of the attribute value.</td></tr><tr><td>pr</td><td>If the attribute has a value, there is a match.</td></tr><tr><td>npr</td><td>If the attribute does not have a value, there is a match.</td></tr><tr><td>gt</td><td>If the attribute value is greater than the operator value, there is a match.  The actual comparison is dependent on the attribute type.</td></tr><tr><td>ge</td><td>If the attribute value is greater than or equal to the operator value, there is a match.  The actual comparison is dependent on the attribute type.</td></tr><tr><td>lt</td><td>If the attribute value is less than the operator value, there is a match.  The actual comparison is dependent on the attribute type.</td></tr><tr><td>le</td><td>If the attribute value is less than or equal to the operator value, there is a match.  The actual comparison is dependent on the attribute type.</td></tr></table><br><br>Example search queries:<table><tr><td>filter=userName eq "bob"&attributes=userName</td></tr><tr><td>filter=name.familyName eq "Marley"&attributes=name</td></tr><tr><td>filter=meta.created ge "2011-09-20T00:00:00Z" and meta.created le "2021-09-21T00:00:00Z"&attributes=userName,meta.created,emails&sortBy=userName&count=2500</td></tr><tr><td>filter=urn:ietf:params:scim:schemas:extension:ibm:2.0:User:customAttributes.favoriteColor eq "blue"&attributes=userName,urn:ietf:params:scim:schemas:extension:ibm:2.0:User:customAttributes.favoriteColor&count=2500</td></tr><tr><td>filter=urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:department eq "2A"&attributes=userName,emails,urn:ietf:params:scim:schemas:extension:enterprise:2.0:User:manager</td></tr><tr><td>filter=emails ew "@example.com" and (phoneNumbers eq "15551212" or phoneNumbers eq "1(555)1212")&attributes=userName,phoneNumbers,emails</td></tr><tr><td>Note: There are some special syntax for "phoneNumbers" to allow filtering using the type, such as GET /v2.0/Users?filter=phoneNumbers.work eq "{value}"&attributes=phoneNumbers.work</td></tr></table><br><br>For tenants that support large groups, additional feature are available.  They are:<table><tr><td>- Search for users in a specific group by using the "memberOf" SCIM attribute.  For example.  GET /v2.0/Users?filter=userName sw "patel" and memberOf eq "{group ID}"</td></tr><tr><td>- Restrict HelpDesk administrators to manage specific groups of users by using Admin Roles.</td></tr></table><br>To check whether the tenant supports large groups, run the GET /v2.0/SCIM/capabilities API.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.GetUsersRequest{}
client.UsersManagementVersion20.GetUsers(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**filter:** `*string` — The SCIM compliant search filter. For example, userName eq "john". The filter should be no longer than 4096 characters in length.
    
</dd>
</dl>

<dl>
<dd>

**attributes:** `*string` — The list of attributes that are passed in as comma-separated values that are used when passing the result back to the caller. To improve performance, specify in the list only the attributes that you want returned. If no list is provided, the default action is to return all attributes.
    
</dd>
</dl>

<dl>
<dd>

**count:** `*string` — Specifies the maximum number of query results per page. A negative value is interpreted as 0.  A value of 0 indicates that no resource results are to be returned, except for totalResults.
    
</dd>
</dl>

<dl>
<dd>

**startIndex:** `*string` — A 1-based index that indicates the start index that is used when the number of users is returned. A value less than 1 is interpreted as 1.
    
</dd>
</dl>

<dl>
<dd>

**sortBy:** `*string` — Sort the results by the specified criteria when the users are returned.
    
</dd>
</dl>

<dl>
<dd>

**sortOrder:** `*generated.GetUsersRequestSortOrder` — The sorting order when the number of users is returned.
    
</dd>
</dl>

<dl>
<dd>

**hashed:** `*string` — The comma-separated list of attributes whose values are to be hashed.
    
</dd>
</dl>

<dl>
<dd>

**fullText:** `*string` — A string that is searched for in the user records.
    
</dd>
</dl>

<dl>
<dd>

**includeGroups:** `*string` — Include group information in the response.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.UsersManagementVersion20.CreateUser(request) -> *generated.UserResponseV2</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlement required: manageUserGroups (Manage users and groups), or manageAllUserGroups (Synchronize users and groups), or manageUserStandardGroups (Manage users and standard groups).<br>Note: You only need one entitlement, but you can have more than one.<br><br>The users are created for a specific tenant that is specified in the request.  Users are either created to use Cloud Directory as an identity source or as a just-in-time provisioning sequence when the user is authenticated at a remote identity source such as an enterprise authentication.<br><br>By default, an email is sent with the password to the user that was created, unless its a federated user. Federated users do not get an email notification. The email templates for branding are at "notifications/user_management/profile/{locale}/account_created_email.xml" and "notifications/user_management/profile/{locale}/account_created_email_with_no_password.xml". Pass in the themeId query parameter to brand the email templates for notifications. To turn off email notifications, send the notifications option  "urn:ietf:params:scim:schemas:extension:ibm:2.0:Notification": {"notifyType":"NONE"} in the payload.<br><br>If custom password intelligence warning is enabled and a password is provided that is listed in it, the 201 response includes the header 'isv-dictionary-policy' with the value: 'WARNLOCAL'.<br>If X-Force password intelligence warning is enabled and a password is provided that is listed in it, the 201 response includes the header 'isv-dictionary-policy' with the value: 'WARNGLOBAL'.<br>If custom password intelligence prevention is enabled and a password is provided that is listed in it, the 400 response can include the header 'isv-dictionary-policy' with the value: 'ENFORCELOCAL'. The corresponding error status is 'PWD_IN_DICTIONARY'.<br>If X-Force password intelligence prevention is enabled and a password is provided that is listed in it, the 400 response can include the header 'isv-dictionary-policy' with the value: 'ENFORCEGLOBAL'. The corresponding error status is 'PWD_IN_GLOBAL_DICTIONARY'.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.CreateUserRequest{
    Body: &generated.UserV2{
        Schemas: []string{
            "schemas",
        },
        UserName: "userName",
    },
}
client.UsersManagementVersion20.CreateUser(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**hashed:** `*string` — The comma-separated list of attributes whose values are to be hashed.
    
</dd>
</dl>

<dl>
<dd>

**themeID:** `*string` — The identifier of the theme that you want to apply.
    
</dd>
</dl>

<dl>
<dd>

**usershouldnotneedtoresetpassword:** `string` — If set to true, the user is not required to change the password after login.<br>Only honored when the password element of UserV2 is set.
    
</dd>
</dl>

<dl>
<dd>

**request:** `*generated.UserV2` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.UsersManagementVersion20.GetUser0(ID) -> *generated.UserResponseV2</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlement required: readUserGroups (Read users and groups), or manageUserGroups (Manage users and groups), or manageAllUserGroups (Synchronize users and groups), or manageUserStandardGroups (Manage users and standard groups).<br>Note: You only need one entitlement, but you can have more than one.<br><br>To improve performance, specify the list of the attributes that you want returned by using the attributes query parameter.<br><br>On Success, the returned response includes the user and group membership details. The group membership that is returned includes the group ID and displayName attributes.<br><br>The memberAttributes, memberCount, and memberStartIndex query parameters are currently ignored.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.GetUser0Request{
    ID: "id",
}
client.UsersManagementVersion20.GetUser0(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The identifier of the user whose details are being retrieved.
    
</dd>
</dl>

<dl>
<dd>

**attributes:** `*string` — The list of attributes that are passed in as comma-separated values that are used when passing the result back to the caller. To improve performance, specify in the list only the attributes that you want returned. If no list is provided, the default action is to return all attributes.
    
</dd>
</dl>

<dl>
<dd>

**memberAttributes:** `*string` — The list of group attributes that are passed in as comma separated values that are used when passing the result back to the caller. For example, the ID and displayName attributes.
    
</dd>
</dl>

<dl>
<dd>

**memberCount:** `*string` — Specifies the maximum number of query results per page. A negative value is interpreted as 0.  A value of 0 indicates that no member resource results are to be returned, except for totalResults.
    
</dd>
</dl>

<dl>
<dd>

**memberStartIndex:** `*string` — The starting index of the search.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.UsersManagementVersion20.PutUser0(ID, request) -> *generated.UserResponseV2</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlement required: manageUserGroups (Manage users and groups), or manageAllUserGroups (Synchronize users and groups), or manageUserStandardGroups (Manage users and standard groups), or updateAnyUser (Update any user).<br>Note: You only need one entitlement, but you can have more than one.<br><br> On Success, the return response contains the user and group membership details. The HTTP PUT method is used to replace the resource's attributes.  For example, clients that previously retrieved the entire resource and revised it, can replace the resource by using an HTTP PUT.<br><br>Only certain attributes for federated users can be modified unless you have the manageAllUserGroups entitlement.  The user's groups cannot be modified and are ignored in the payload. Group membership is managed by using the PUT /v2.0/Groups/{id} and PATCH /v2.0/Groups/{id} API.<br><br>By default, an email is sent to regular users that includes the changed attributes. Federated users do not get an email notification. The email template for branding is at "notifications/user_management/profile/{locale}/user_profile_modified_email.xml". Pass in the themeId query parameter to brand the email template for notifications. To turn off email notifications, send the notifications option  "urn:ietf:params:scim:schemas:extension:ibm:2.0:Notification": {"notifyType":"NONE"} in the payload.<br><br>If custom password intelligence warning is enabled and a password is provided that is listed in it, the 200 response includes the header 'isv-dictionary-policy' with the value: 'WARNLOCAL'.<br>If X-Force password intelligence warning is enabled and a password is provided that is listed in it, the 200 response includes the header 'isv-dictionary-policy' with the value: 'WARNGLOBAL'.<br>If custom password intelligence prevention is enabled and a password is provided that is listed in it, the 400 response can include the header 'isv-dictionary-policy' with the value: 'ENFORCELOCAL'. The corresponding error status is 'PWD_IN_DICTIONARY'.<br>If X-Force password intelligence prevention is enabled and a password is provided that is listed in it, the 400 response can include the header 'isv-dictionary-policy' with the value: 'ENFORCEGLOBAL'. The corresponding error status is 'PWD_IN_GLOBAL_DICTIONARY'.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.PutUser0Request{
    ID: "id",
    Body: &generated.UserV2{
        Schemas: []string{
            "schemas",
        },
        UserName: "userName",
    },
}
client.UsersManagementVersion20.PutUser0(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` 
    
</dd>
</dl>

<dl>
<dd>

**hashed:** `*string` — The comma separated list of attributes whose values are to be hashed.
    
</dd>
</dl>

<dl>
<dd>

**themeID:** `*string` — The identifier of the theme that you want to apply.
    
</dd>
</dl>

<dl>
<dd>

**usershouldnotneedtoresetpassword:** `string` — If set to true for a password change, the user is not required to change the password after login.<br>Only honored when the password element of UserV2 is set.
    
</dd>
</dl>

<dl>
<dd>

**request:** `*generated.UserV2` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.UsersManagementVersion20.DeleteUser0(ID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlement required: manageUserGroups (Manage users and groups), or manageAllUserGroups (Synchronize users and groups), or manageUserStandardGroups (Manage users and standard groups).<br>Note: You only need one entitlement, but you can have more than one.<br><br>By default, an email is sent to regular users that the account was deleted. Federated users do not get an email notification. The email template for branding is at "notifications/user_management/profile/{locale}/account_deleted_email.xml". Pass in the themeId query parameter to brand the email template for notifications.<br><br>To turn off email notifications, send notifyType=NONE as a query parameter.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.DeleteUser0Request{
    ID: "id",
}
client.UsersManagementVersion20.DeleteUser0(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The identifier of the user that is being deleted.
    
</dd>
</dl>

<dl>
<dd>

**notifyType:** `*generated.DeleteUser0RequestNotifyType` — An optional query parameter that denotes the notification type.  If not present, the EMAIL notification is used. Specify NONE if no notification to the user that their account has been deleted is required.
    
</dd>
</dl>

<dl>
<dd>

**themeID:** `*string` — The identifier of the theme that you want to apply.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.UsersManagementVersion20.PatchUser(ID, request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlement required: manageUserGroups (Manage users and groups), or manageAllUserGroups (Synchronize users and groups), or manageUserStandardGroups (Manage users and standard groups), or updateAnyUser (Update any user). <br>Note: You only need one entitlement, but you can have more than one.<br><br>Only certain attributes for federated users can be modified unless you have the manageAllUserGroups entitlement.  The user's groups cannot be modified and are ignored in the payload. Group membership is managed by using the PUT /v2.0/Groups/{id} and PATCH /v2.0/Groups/{id} API.<br><br>By default, an email is sent to regular users that includes the changed attributes. Federated users do not get an email notification. The email template for branding is at "notifications/user_management/profile/{locale}/user_profile_modified_email.xml". Pass in the themeId query parameter to brand the email template for notifications. <br><br>The following is an example of a patch request that adds a title, replaces the formatted name, and removes any custom attributes which name contains "customA" from the user. Notice also that the example shows how you can specify the notifyType if you want to.<br>NotifyType is an optional attribute that denotes the notification type.  If not present, the EMAIL notification is used.Specify NONE if no notification is required.<br><br> `{ "schemas":["urn:ietf:params:scim:api:messages:2.0:PatchOp"],"Operations":  [    {      "op":"add",      "path":"title",      "value":"Vice President"    },    {      "op":"replace",      "path":"name.formatted",      "value":"James Smith-Miller"    },    {      "op":"remove",      "path":"urn:ietf:params:scim:schemas:extension:ibm:2.0:User:customAttributes[name co \"customA\"]"    },    {      "op": "add",      "path": "urn:ietf:params:scim:schemas:extension:ibm:2.0:Notification:notifyType",      "value": "EMAIL"    }  ]}` <br><br>If custom password intelligence warning is enabled and a password is provided that is listed in it, the 204 response includes the header 'isv-dictionary-policy' with the value: 'WARNLOCAL'.<br>If X-Force password intelligence warning is enabled and a password is provided that is listed in it, the 204 response includes the header 'isv-dictionary-policy' with the value: 'WARNGLOBAL'.<br>If custom password intelligence prevention is enabled and a password is provided that is listed in it, the 400 response can include the header 'isv-dictionary-policy' with the value: 'ENFORCELOCAL'. The corresponding error status is 'PWD_IN_DICTIONARY'.<br>If X-Force password intelligence prevention is enabled and a password is provided that is listed in it, the 400 response can include the header 'isv-dictionary-policy' with the value: 'ENFORCEGLOBAL'. The corresponding error status is 'PWD_IN_GLOBAL_DICTIONARY'.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.PatchBody{
    ID: "id",
    Schemas: []string{
        "schemas",
    },
    Operations: []*generated.PatchOperation0{
        &generated.PatchOperation0{
            Op: generated.PatchOperation0OpAdd,
            Path: "title",
        },
    },
}
client.UsersManagementVersion20.PatchUser(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**id:** `string` — The ID of the user to be patched.
    
</dd>
</dl>

<dl>
<dd>

**themeID:** `*string` — The identifier of the theme that you want to apply.
    
</dd>
</dl>

<dl>
<dd>

**usershouldnotneedtoresetpassword:** `string` — If set to true for a password change, the user is not required to change the password after login.
    
</dd>
</dl>

<dl>
<dd>

**schemas:** `[]string` — The body of each SCIM PATCH request must contain the "schemas" attribute with the URI value: "urn:ietf:params:scim:api:messages:2.0:PatchOp".
    
</dd>
</dl>

<dl>
<dd>

**operations:** `[]*generated.PatchOperation0` — An array of operation objects to be performed.  Operation objects must have exactly one "op" member, whose value indicates the operation to perform. Its value must be one of "add", "remove", or "replace". Values are errors.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Application Access
<details><summary><code>client.ApplicationAccess.GetApplication(ApplicationID) -> *generated.ApplicationDetailsResponseBean</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlements required:  manageAppAccessAdmin  (Manage application lifecycle)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.GetApplicationRequest{
    ApplicationID: "applicationId",
}
client.ApplicationAccess.GetApplication(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**applicationID:** `string` — The application ID for which details are being requested.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.ApplicationAccess.UpdateApplication(ApplicationID, request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Rest interface to update the specified application. Entitlements required: manageAppAccessAdmin  (Manage application lifecycle)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.UpdateApplicationRequest{
    ApplicationID: "applicationId",
    Body: &generated.ApplicationRequestBean{
        Name: "name",
        TemplateID: "templateId",
        Providers: &generated.ProviderBean{
            SSO: &generated.SSOBean{
                DomainName: "domainName",
            },
        },
        Provisioning: &generated.ProvisioningBean{},
    },
}
client.ApplicationAccess.UpdateApplication(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**applicationID:** `string` — ID of the application to be updated.
    
</dd>
</dl>

<dl>
<dd>

**request:** `*generated.ApplicationRequestBean` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.ApplicationAccess.DeleteApplication(ApplicationID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlements required:  manageAppAccessAdmin  (Manage application lifecycle)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.DeleteApplicationRequest{
    ApplicationID: "applicationId",
}
client.ApplicationAccess.DeleteApplication(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**applicationID:** `string` — The application ID for which deletion requested.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.ApplicationAccess.SearchApplications() -> *generated.SearchAdminApplicationWithoutProvResponseBean</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlements required:  manageAppAccessAdmin  (Manage application lifecycle). Under the provisioning section the attributeMappings, reverseAttributeMappings and extension sections will always be an empty object.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.SearchApplicationsRequest{}
client.ApplicationAccess.SearchApplications(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**page:** `*string` — Ordinality of the page to fetch. Required for paginated queries.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*string` — Total number of applications to return in each page. Cannot be greater than count. Required for paginated queries.
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — Search query to return matching applications e.g search="q={searchString}". The query returns list of application which has {searchString} as a sub-string in application name. Examples: 1. If a tenant has Applications with name App1, App12, App123 and Query is search="q=pp12" then the applications returned will be App12 & App123. If search="q=pp1", then all the applications which contains App1 in it will be returned i.e. App1, App12, App123.
    
</dd>
</dl>

<dl>
<dd>

**sort:** `*string` — Attributes to sort results on, supported values are 'name' and 'entityid'. Prepend the attribute with '+' or '-' sign for ascending and descending sorted order respectively. If not specified, sorted in ascending order on entityid. The entity id corresponds to application id.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.ApplicationAccess.CreateApplication(request) -> *generated.PostApplicationResponseBean</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Entitlements required:  manageAppAccessAdmin  (Manage application lifecycle). If Application is enabled for Account Lifecycle or Account Sync, ensure the credentials specified for target endpoint are correct. This can be verified by performing Test Connection in Account Lifecycle section of Application settings.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.ApplicationRequestBean{
    Name: "name",
    TemplateID: "templateId",
    Providers: &generated.ProviderBean{
        SSO: &generated.SSOBean{
            DomainName: "domainName",
        },
    },
    Provisioning: &generated.ProvisioningBean{},
}
client.ApplicationAccess.CreateApplication(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**request:** `*generated.ApplicationRequestBean` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## API Clients
<details><summary><code>client.APIClients.GetAPIClient(ClientID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Gets the API client with the specified ID. Note: The response contains the API client ID and secret. The API client secret is privileged information.<br><br>Entitlements required: manageAPIClients (Manage API clients) or readAPIClients (Read API clients)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.GetAPIClientRequest{
    ClientID: "clientId",
}
client.APIClients.GetAPIClient(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**clientID:** `string` — The unique ID of a client.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.APIClients.UpdateAPIClient(ClientIDPathParam, request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates the specified API client properties except for the ID and the client ID which cannot be modified. If the client secret is not specified, then a new client secret will be generated.<br><br>The entitlements array can contain any combination of entitlements.<br><br><details><summary>List of API Client entitlements:</summary><br><table><tr><th>Entitlement</th><th>Description</th><th>Offering</th></tr><tr><td>manageDeployment</td><td>Manage deployment</td><td>any</td></tr><tr><td>manageCerts</td><td>Manage certificates</td><td>any</td></tr><tr><td>readCerts</td><td>Read certificates</td><td>any</td></tr><tr><td>manageAPIClients</td><td>Manage API clients</td><td>any</td></tr><tr><td>readAPIClients</td><td>Read API clients</td><td>any</td></tr><tr><td>manageIdentitySources</td><td>Manage identity providers</td><td>any</td></tr><tr><td>readIdentitySources</td><td>Read identity providers</td><td>any</td></tr><tr><td>manageMFAMethods</td><td>Manage second-factor authentication method configuration</td><td>CIC</td></tr><tr><td>readMFAMethods</td><td>Read second-factor authentication method configuration</td><td>CIC</td></tr><tr><td>manageEnrollMFAMethodAnyUser</td><td>Manage second-factor authentication enrollment for all users</td><td>CIV</td></tr><tr><td>readEnrollMFAMethodAnyUser</td><td>Read second-factor authentication enrollment for all users</td><td>CIV</td></tr><tr><td>authnAnyUser</td><td>Authenticate any user</td><td>CIV</td></tr><tr><td>manageAuthenticatorsConfig</td><td>Manage authenticator configuration</td><td>CIV</td></tr><tr><td>readAuthenticatorsConfig</td><td>Read authenticator configuration</td><td>CIV</td></tr><tr><td>manageAuthenticatorsAnyUser</td><td>Manage authenticator registrations for all users</td><td>CIV</td></tr><tr><td>readAuthenticatorsAnyUser</td><td>Read authenticator registrations for all users</td><td>CIV</td></tr><tr><td>manageUserGroups</td><td>Manage users and groups</td><td>any</td></tr><tr><td>readUserGroups</td><td>Read users and groups</td><td>any</td></tr><tr><td>manageAllUserGroups</td><td>Synchronize users and groups</td><td>any</td></tr><tr><td>manageUsersPwdReset</td><td>Manage users and their pwdReset attribute</td><td>any</td></tr><tr><td>manageUserStandardGroups</td><td>Manage users and standard groups</td><td>any</td></tr><tr><td>manageAdminGroup</td><td>Manage administrator group</td><td>any</td></tr><tr><td>readAdminGroup</td><td>Read administrator group</td><td>any</td></tr><tr><td>managePwdPolicy</td><td>Manage password policy</td><td>any</td></tr><tr><td>readPwdPolicy</td><td>Read password policy</td><td>any</td></tr><tr><td>AnalyticsDataSyncToCloud</td><td></td><td>CIA</td></tr><tr><td>AnalyticsSatelliteOnBoard</td><td></td><td>CIA</td></tr><tr><td>manageOIDCGrants</td><td>Manage OAuth tokens</td><td>any</td></tr><tr><td>readOIDCGrants</td><td>Read OAuth tokens</td><td>any</td></tr><tr><td>recoverUsername</td><td>Recover user name</td><td>any</td></tr><tr><td>manageFederations</td><td>Manage federations	</td><td>any</td></tr><tr><td>readFederations</td><td>Read federations	</td><td>any</td></tr><tr><td>resetPassword</td><td>Reset password	</td><td>any</td></tr><tr><td>manageAppAccessAdmin</td><td>Manage application lifecycle</td><td>any</td></tr><tr><td>manageAppAccessOwner</td><td>Manage application entitlements</td><td>any</td></tr><tr><td>manageSubscriptions</td><td>Manage subscriptions</td><td>ISC</td></tr><tr><td>manageAccessPolicies</td><td>Manage access policies</td><td>any</td></tr><tr><td>readAccessPolicies</td><td>Read access policies</td><td>any</td></tr><tr><td>managePushCreds</td><td>Manage Push notification credentials</td><td>any</td></tr><tr><td>readPushCreds</td><td>Read Push notification credentials</td><td>any</td></tr><tr><td>manageAccessRequest</td><td>Manage access request</td><td>CIG</td></tr><tr><td>manageAccessWorkflow</td><td>Manage access request work flows</td><td>CIG</td></tr><tr><td>manageOIDCConsents</td><td>Manage OAuth consents</td><td>any</td></tr><tr><td>readOIDCConsents</td><td>Read OAuth consents</td><td>any</td></tr><tr><td>manageReports</td><td>Manage reports</td><td>any. Exception: application usage reports can only be exported by CIC.</td></tr><tr><td>readReports</td><td>Read reports</td><td>any. Exception: application usage reports can only be accessed by CIC.</td></tr><tr><td>updateAnyUser</td><td>Update any user</td><td>any</td></tr><tr><td>resetPasswordAnyUser</td><td>Reset password of any user</td><td>any</td></tr><tr><td>readTenantProperties</td><td>Read tenant properties</td><td>any</td></tr><tr><td>manageTenantProperties</td><td>Manage tenant properties</td><td>any</td></tr><tr><td>manageAttributes</td><td>Manage attribute sources</td><td>any</td></tr><tr><td>readAttributes</td><td>Read attribute sources</td><td>any</td></tr><tr><td>generateOTP</td><td>Generate OTP</td><td>CIV</td></tr><tr><td>readAppConfig</td><td>Read application configuration</td><td>any</td></tr><tr><td>manageTemplates</td><td>Manage templates and themes</td><td>any</td></tr><tr><td>readTemplates</td><td>Read templates and themes</td><td>any</td></tr><tr><td>reviewCertRecords</td><td>Review certification records</td><td>CIG</td></tr><tr><td>readEntitlements</td><td>Read configurable entitlements</td><td>any</td></tr><tr><td>manageNotificationProviders</td><td>Manage notification providers</td><td>any</td></tr><tr><td>readNotificationProviders</td><td>Read notification providers</td><td>any</td></tr><tr><td>manageCertifications</td><td>Manage certifications</td><td>CIG</td></tr><tr><td>readExternalAgents</td><td>Read external agents</td><td>any</td></tr><tr><td>manageExternalAgents</td><td>Manage external agents</td><td>any</td></tr><tr><td>runExternalAgent</td><td>Enable external agent runtime functions</td><td>any</td></tr><tr><td>manageOidcDynamicClient</td><td>Manage OIDC client registration dynamically</td><td>any</td></tr><tr><td>readPurpose</td><td>Read privacy purposes and EULA</td><td>any</td></tr><tr><td>managePurpose</td><td>Manage privacy purposes and EULA</td><td>any</td></tr><tr><td>manageAppPurpose</td><td>Manage application privacy purposes</td><td>any</td></tr><tr><td>readPrivacyConsent</td><td>Read privacy consents</td><td>any</td></tr><tr><td>managePrivacyConsent</td><td>Manage privacy consents</td><td>any</td></tr><tr><td>readPrivacyPolicy</td><td>Read privacy rules and policy</td><td>any</td></tr><tr><td>managePrivacyPolicy</td><td>Manage privacy rules and policy</td><td>any</td></tr><tr><td>createPrivacyConsent</td><td>Create privacy consent records</td><td>any</td></tr><tr><td>performDSP</td><td>Retrieve privacy purposes and associated user's consent</td><td>any</td></tr><tr><td>performDUA</td><td>Check for data usage approval</td><td>any</td></tr><tr><td>certCampaignSupervisor</td><td>Monitor certification campaigns</td><td>CIG</td></tr><tr><td>managePwdVaultAnyUser</td><td>Manage password vault for all users</td><td>CIC, CIV</td></tr><tr><td>managePwdVault</td><td>Manage own password vault</td><td>CIC, CIV</td></tr><tr><td>readPwdVaultAnyUser</td><td>Read password vault for all users</td><td>CIC, CIV</td></tr><tr><td>readPwdVault</td><td>Read own password vault</td><td>CIC, CIV</td></tr><tr><td>managePwdVaultConfig</td><td>Manage password vault configuration</td><td>CIC, CIV</td></tr><tr><td>readPwdVaultConfig</td><td>Read password vault configuration</td><td>CIC, CIV</td></tr><tr><td>mfaPush</td><td>Send second-factor push notifications</td><td>CIV</td></tr><tr><td>readPrivacyProfile</td><td>Read privacy profiles</td><td>any</td></tr><tr><td>managePrivacyProfile</td><td>Manage privacy profiles</td><td>any</td></tr><tr><td>manageEntitlements</td><td>Manage entitlements</td><td>any</td></tr><tr><td>manageDevicesAnyUser</td><td>Manage devices for all users</td><td>any</td></tr><tr><td>readDevicesAnyUser</td><td>Read devices for all users</td><td>any</td></tr><tr><td>manageDevices</td><td>Manage only your devices</td><td>any</td></tr><tr><td>readDevices</td><td>Read only your devices</td><td>any</td></tr><tr><td>manageRecaptcha</td><td>Manage reCAPTCHA configuration</td><td>any</td></tr><tr><td>readRecaptcha</td><td>Read reCAPTCHA configuration</td><td>any</td></tr><tr><td>manageLoginSessions</td><td>Manage login sessions</td><td>any</td></tr><tr><td>manageRelyingParty</td><td>Manage relying party configuration</td><td>any</td></tr><tr><td>readRelyingParty</td><td>Read relying party configuration</td><td>any</td></tr><tr><td>manageWebhooks</td><td>Manage webhooks <td>any</td></tr><tr><td>readWebhooks</td><td>Read webhooks </td><td>any</td></tr><tr><td>readSTSClients</td><td>Read STS clients and token types</td><td>any</td></tr><tr><td>manageSTSClients</td><td>Manage STS clients and token types</td><td>any</td></tr><tr><td>manageVerifiableLinks</td><td>Manage verifiable links configuration</td><td>any</td></tr><tr><td>readSelfOidcGrants</td><td>Read your OIDC and OAuth grants</td><td>any</td></tr><tr><td>manageSelfOidcGrants</td><td>Manage your OIDC and OAuth grants</td><td>any</td></tr><tr><td>diManageAgency</td><td>Manage Decentralized Identity Agency Configuration</td><td>any</td></tr><tr><td>diReadAgency</td><td>Read Decentralized Identity Agency Configuration</td><td>any</td></tr><tr><td>diManageAgentsAny</td><td>Manage Decentralized Identity Agents</td><td>any</td></tr><tr><td>diReadAgentsAny</td><td>Read Decentralized Identity Agents</td><td>any</td></tr><tr><td>manageMyOrg</td><td>Manage my organization</td><td>CIG</td></tr><tr><td>diIssueCredentials</td><td>Issue Decentralized Identity Verifiable Credentials</td><td>CIV</td></tr><tr><td>diVerifyCredentials</td><td>Verify Decentralized Identity Verifiable Credentials</td><td>CIV</td></tr></table></details><br><br>Entitlements required: manageAPIClients (Manage API clients)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.APIClientConfig{
    ClientIDPathParam: "clientId",
    ClientName: "Demo Client",
    Entitlements: []string{
        "entitlements",
    },
}
client.APIClients.UpdateAPIClient(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**clientIDPathParam:** `string` — The unique ID of a client.
    
</dd>
</dl>

<dl>
<dd>

**clientName:** `string` — the friendly name of the client
    
</dd>
</dl>

<dl>
<dd>

**entitlements:** `[]string` — the list of entitlements assigned to the client
    
</dd>
</dl>

<dl>
<dd>

**clientID:** `*string` — the generated client id for authorization
    
</dd>
</dl>

<dl>
<dd>

**clientSecret:** `*string` — the generated client secret for authorization
    
</dd>
</dl>

<dl>
<dd>

**enabled:** `*bool` — whether or not the client can be used to generate tokens
    
</dd>
</dl>

<dl>
<dd>

**overrideSettings:** `*generated.APIClientOverrideSettings` 
    
</dd>
</dl>

<dl>
<dd>

**description:** `*string` — a description of the client
    
</dd>
</dl>

<dl>
<dd>

**additionalProperties:** `map[string]*generated.APIClientConfigAdditionalPropertiesValue` — additional properties for the client
    
</dd>
</dl>

<dl>
<dd>

**ipFilterOp:** `*generated.APIClientConfigIPFilterOp` — the operation of the ip filter. The default setting is null, which means that the ip filter is disabled
    
</dd>
</dl>

<dl>
<dd>

**ipFilters:** `[]string` — the list of ips
    
</dd>
</dl>

<dl>
<dd>

**jwkURI:** `*string` — the JSON web key URI endpoint
    
</dd>
</dl>

<dl>
<dd>

**additionalConfig:** `*generated.APIClientAdditionalConfig` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.APIClients.DeleteAPIClient(ClientID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes the specified API client. A successful response is returned even if the object does not exist.<br><br>Entitlements required: manageAPIClients (Manage API clients)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.DeleteAPIClientRequest{
    ClientID: "clientId",
}
client.APIClients.DeleteAPIClient(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**clientID:** `string` — The unique ID of a client
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.APIClients.GetAPIClients() -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Gets a list of the configured API clients. Note: The response contains the API client ID and secret for each API client. The API client secret is privileged information.<br><br>Entitlements required: manageAPIClients (Manage API clients) or readAPIClients (Read API clients)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.GetAPIClientsRequest{}
client.APIClients.GetAPIClients(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**pagination:** `*string` — The prefix for the paging parameter is "pagination=". If no pagination parameters are passed in, all results are returned. The maximum allowed value for limit or count is 1000. <br><b>count</b> is the total number of results to be returned from the data store.<br><b>page</b> is which page we are requesting, or the offset. <br><b>limit</b> is the total number of results to return in one page.<br><br>The pagination parameter value <b>must</b> be HTML encoded.<br>Note: This is not required when using the Swagger UI.<br><br><b>Example:</b> Paginate on count=10&page=1&limit=5<br>pagination=count%3D10%26page%3D1%26limit%3D5
    
</dd>
</dl>

<dl>
<dd>

**sort:** `*string` — The prefix for the sort parameter is "sort=". Each attribute must be prefixed with either + or - (+ ascending, - descending). Multiple attributes must be separated by a comma (,).<br><br>The valid fields for sorting are: clientId, clientName, and enabled.<br><br>The sort parameter value <b>must</b> be HTML encoded.<br>Note: This is not required when using the Swagger UI.<br><br><b>Example:</b> Sort on -enabled,+clientId<br>sort=-enabled%2C%2BclientId
    
</dd>
</dl>

<dl>
<dd>

**search:** `*string` — The prefix for all search operations is "search=".<br>Valid operators for strings are = , !=  and contains <br>Valid operators for booleans are = and !=<br>Valid operators for numbers are >=, >, &lt=, &lt, = and !=<br>String search values must be double quoted, numbers and booleans must not.<br><br>The valid fields for sorting are: clientId, clientName, and enabled.<br><br>The search parameter value <b>must</b> be HTML encoded.<br>Note: This is not required when using the Swagger UI.<br><br><b>Example:</b> Search on clientId contains "ABCDEF"&enabled=true<br>search=clientId%20contains%20%22ABCDEF%22%26enabled%3Dtrue
    
</dd>
</dl>

<dl>
<dd>

**filter:** `*string` — The prefix for the filter parameter is "filter="<br>Valid formats are either inclusive only or exclusive only. These must not be intermingled. Multiple filter parameters must be separated by a comma (,).<br><br>The valid fields for filtering are: id, clientId, clientName, clientSecret, entitlements, and enabled.<br><br>The filter parameter value <b>must</b> be HTML encoded.<br>Note: This is not required when using the Swagger UI.<br><br><b>Examples</b><br>Filter to only return clientId:<br>filter=clientId<br><br>Filter to exclude clientSecret and enabled:<br>filter=%21clientSecret,enabled
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.APIClients.CreateAPIClient(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates an API client with a random client ID and secret and assigns the given entitlements. The client is configured with the client_credentials grant type. You must perform a GET operation on the returned location header to get the generated client ID and secret.<br><br>The entitlements array can contain any combination of entitlements.<br><br><details><summary>List of API Client entitlements:</summary><br><table><tr><th>Entitlement</th><th>Description</th><th>Offering</th></tr><tr><td>manageDeployment</td><td>Manage deployment</td><td>any</td></tr><tr><td>manageCerts</td><td>Manage certificates</td><td>any</td></tr><tr><td>readCerts</td><td>Read certificates</td><td>any</td></tr><tr><td>manageAPIClients</td><td>Manage API clients</td><td>any</td></tr><tr><td>readAPIClients</td><td>Read API clients</td><td>any</td></tr><tr><td>manageIdentitySources</td><td>Manage identity providers</td><td>any</td></tr><tr><td>readIdentitySources</td><td>Read identity providers</td><td>any</td></tr><tr><td>manageMFAMethods</td><td>Manage second-factor authentication method configuration</td><td>CIC</td></tr><tr><td>readMFAMethods</td><td>Read second-factor authentication method configuration</td><td>CIC</td></tr><tr><td>manageEnrollMFAMethodAnyUser</td><td>Manage second-factor authentication enrollment for all users</td><td>CIV</td></tr><tr><td>readEnrollMFAMethodAnyUser</td><td>Read second-factor authentication enrollment for all users</td><td>CIV</td></tr><tr><td>authnAnyUser</td><td>Authenticate any user</td><td>CIV</td></tr><tr><td>manageAuthenticatorsConfig</td><td>Manage authenticator configuration</td><td>CIV</td></tr><tr><td>readAuthenticatorsConfig</td><td>Read authenticator configuration</td><td>CIV</td></tr><tr><td>manageAuthenticatorsAnyUser</td><td>Manage authenticator registrations for all users</td><td>CIV</td></tr><tr><td>readAuthenticatorsAnyUser</td><td>Read authenticator registrations for all users</td><td>CIV</td></tr><tr><td>manageUserGroups</td><td>Manage users and groups</td><td>any</td></tr><tr><td>readUserGroups</td><td>Read users and groups</td><td>any</td></tr><tr><td>manageAllUserGroups</td><td>Synchronize users and groups</td><td>any</td></tr><tr><td>manageUsersPwdReset</td><td>Manage users and their pwdReset attribute</td><td>any</td></tr><tr><td>manageUserStandardGroups</td><td>Manage users and standard groups</td><td>any</td></tr><tr><td>manageAdminGroup</td><td>Manage administrator group</td><td>any</td></tr><tr><td>readAdminGroup</td><td>Read administrator group</td><td>any</td></tr><tr><td>managePwdPolicy</td><td>Manage password policy</td><td>any</td></tr><tr><td>readPwdPolicy</td><td>Read password policy</td><td>any</td></tr><tr><td>AnalyticsDataSyncToCloud</td><td></td><td>CIA</td></tr><tr><td>AnalyticsSatelliteOnBoard</td><td></td><td>CIA</td></tr><tr><td>manageOIDCGrants</td><td>Manage OAuth tokens</td><td>any</td></tr><tr><td>readOIDCGrants</td><td>Read OAuth tokens</td><td>any</td></tr><tr><td>recoverUsername</td><td>Recover user name</td><td>any</td></tr><tr><td>manageFederations</td><td>Manage federations	</td><td>any</td></tr><tr><td>readFederations</td><td>Read federations	</td><td>any</td></tr><tr><td>resetPassword</td><td>Reset password	</td><td>any</td></tr><tr><td>manageAppAccessAdmin</td><td>Manage application lifecycle</td><td>any</td></tr><tr><td>manageAppAccessOwner</td><td>Manage application entitlements</td><td>any</td></tr><tr><td>manageSubscriptions</td><td>Manage subscriptions</td><td>ISC</td></tr><tr><td>manageAccessPolicies</td><td>Manage access policies</td><td>any</td></tr><tr><td>readAccessPolicies</td><td>Read access policies</td><td>any</td></tr><tr><td>managePushCreds</td><td>Manage Push notification credentials</td><td>any</td></tr><tr><td>readPushCreds</td><td>Read Push notification credentials</td><td>any</td></tr><tr><td>manageAccessRequest</td><td>Manage access request</td><td>CIG</td></tr><tr><td>manageAccessWorkflow</td><td>Manage access request work flows</td><td>CIG</td></tr><tr><td>manageOIDCConsents</td><td>Manage OAuth consents</td><td>any</td></tr><tr><td>readOIDCConsents</td><td>Read OAuth consents</td><td>any</td></tr><tr><td>manageReports</td><td>Manage reports</td><td>any. Exception: application usage reports can only be exported by CIC.</td></tr><tr><td>readReports</td><td>Read reports</td><td>any. Exception: application usage reports can only be accessed by CIC.</td></tr><tr><td>updateAnyUser</td><td>Update any user</td><td>any</td></tr><tr><td>resetPasswordAnyUser</td><td>Reset password of any user</td><td>any</td></tr><tr><td>readTenantProperties</td><td>Read tenant properties</td><td>any</td></tr><tr><td>manageTenantProperties</td><td>Manage tenant properties</td><td>any</td></tr><tr><td>manageAttributes</td><td>Manage attribute sources</td><td>any</td></tr><tr><td>readAttributes</td><td>Read attribute sources</td><td>any</td></tr><tr><td>generateOTP</td><td>Generate OTP</td><td>CIV</td></tr><tr><td>readAppConfig</td><td>Read application configuration</td><td>any</td></tr><tr><td>manageTemplates</td><td>Manage templates and themes</td><td>any</td></tr><tr><td>readTemplates</td><td>Read templates and themes</td><td>any</td></tr><tr><td>reviewCertRecords</td><td>Review certification records</td><td>CIG</td></tr><tr><td>readEntitlements</td><td>Read configurable entitlements</td><td>any</td></tr><tr><td>manageNotificationProviders</td><td>Manage notification providers</td><td>any</td></tr><tr><td>readNotificationProviders</td><td>Read notification providers</td><td>any</td></tr><tr><td>manageCertifications</td><td>Manage certifications</td><td>CIG</td></tr><tr><td>readExternalAgents</td><td>Read external agents</td><td>any</td></tr><tr><td>manageExternalAgents</td><td>Manage external agents</td><td>any</td></tr><tr><td>runExternalAgent</td><td>Enable external agent runtime functions</td><td>any</td></tr><tr><td>manageOidcDynamicClient</td><td>Manage OIDC client registration dynamically</td><td>any</td></tr><tr><td>readPurpose</td><td>Read privacy purposes and EULA</td><td>any</td></tr><tr><td>managePurpose</td><td>Manage privacy purposes and EULA</td><td>any</td></tr><tr><td>manageAppPurpose</td><td>Manage application privacy purposes</td><td>any</td></tr><tr><td>readPrivacyConsent</td><td>Read privacy consents</td><td>any</td></tr><tr><td>managePrivacyConsent</td><td>Manage privacy consents</td><td>any</td></tr><tr><td>readPrivacyPolicy</td><td>Read privacy rules and policy</td><td>any</td></tr><tr><td>managePrivacyPolicy</td><td>Manage privacy rules and policy</td><td>any</td></tr><tr><td>createPrivacyConsent</td><td>Create privacy consent records</td><td>any</td></tr><tr><td>performDSP</td><td>Retrieve privacy purposes and associated user's consent</td><td>any</td></tr><tr><td>performDUA</td><td>Check for data usage approval</td><td>any</td></tr><tr><td>certCampaignSupervisor</td><td>Monitor certification campaigns</td><td>CIG</td></tr><tr><td>managePwdVaultAnyUser</td><td>Manage password vault for all users</td><td>CIC, CIV</td></tr><tr><td>managePwdVault</td><td>Manage own password vault</td><td>CIC, CIV</td></tr><tr><td>readPwdVaultAnyUser</td><td>Read password vault for all users</td><td>CIC, CIV</td></tr><tr><td>readPwdVault</td><td>Read own password vault</td><td>CIC, CIV</td></tr><tr><td>managePwdVaultConfig</td><td>Manage password vault configuration</td><td>CIC, CIV</td></tr><tr><td>readPwdVaultConfig</td><td>Read password vault configuration</td><td>CIC, CIV</td></tr><tr><td>mfaPush</td><td>Send second-factor push notifications</td><td>CIV</td></tr><tr><td>readPrivacyProfile</td><td>Read privacy profiles</td><td>any</td></tr><tr><td>managePrivacyProfile</td><td>Manage privacy profiles</td><td>any</td></tr><tr><td>manageEntitlements</td><td>Manage entitlements</td><td>any</td></tr><tr><td>manageDevicesAnyUser</td><td>Manage devices for all users</td><td>any</td></tr><tr><td>readDevicesAnyUser</td><td>Read devices for all users</td><td>any</td></tr><tr><td>manageDevices</td><td>Manage only your devices</td><td>any</td></tr><tr><td>readDevices</td><td>Read only your devices</td><td>any</td></tr><tr><td>manageRecaptcha</td><td>Manage reCAPTCHA configuration</td><td>any</td></tr><tr><td>readRecaptcha</td><td>Read reCAPTCHA configuration</td><td>any</td></tr><tr><td>manageLoginSessions</td><td>Manage login sessions</td><td>any</td></tr><tr><td>manageRelyingParty</td><td>Manage relying party configuration</td><td>any</td></tr><tr><td>readRelyingParty</td><td>Read relying party configuration</td><td>any</td></tr><tr><td>manageWebhooks</td><td>Manage webhooks <td>any</td></tr><tr><td>readWebhooks</td><td>Read webhooks </td><td>any</td></tr><tr><td>readSTSClients</td><td>Read STS clients and token types</td><td>any</td></tr><tr><td>manageSTSClients</td><td>Manage STS clients and token types</td><td>any</td></tr><tr><td>manageVerifiableLinks</td><td>Manage verifiable links configuration</td><td>any</td></tr><tr><td>readSelfOidcGrants</td><td>Read your OIDC and OAuth grants</td><td>any</td></tr><tr><td>manageSelfOidcGrants</td><td>Manage your OIDC and OAuth grants</td><td>any</td></tr><tr><td>diManageAgency</td><td>Manage Decentralized Identity Agency Configuration</td><td>any</td></tr><tr><td>diReadAgency</td><td>Read Decentralized Identity Agency Configuration</td><td>any</td></tr><tr><td>diManageAgentsAny</td><td>Manage Decentralized Identity Agents</td><td>any</td></tr><tr><td>diReadAgentsAny</td><td>Read Decentralized Identity Agents</td><td>any</td></tr><tr><td>manageMyOrg</td><td>Manage my organization</td><td>CIG</td></tr><tr><td>diIssueCredentials</td><td>Issue Decentralized Identity Verifiable Credentials</td><td>CIV</td></tr><tr><td>diVerifyCredentials</td><td>Verify Decentralized Identity Verifiable Credentials</td><td>CIV</td></tr></table></details><br><br>Entitlements required: manageAPIClients (Manage API clients)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &generated.APIClientConfigRequest{
    ClientName: "Demo Client",
    Entitlements: []string{
        "entitlements",
    },
    Enabled: true,
}
client.APIClients.CreateAPIClient(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**clientID:** `*string` — the unique identifier of the client
    
</dd>
</dl>

<dl>
<dd>

**clientName:** `string` — the friendly name of the client
    
</dd>
</dl>

<dl>
<dd>

**clientSecret:** `*string` — the generated client secret for authorization. If unspecified, a random client secret is generated
    
</dd>
</dl>

<dl>
<dd>

**entitlements:** `[]string` — the list of entitlements assigned to the client
    
</dd>
</dl>

<dl>
<dd>

**enabled:** `bool` — whether or not the client can be used to generate tokens
    
</dd>
</dl>

<dl>
<dd>

**overrideSettings:** `*generated.APIClientOverrideSettings` 
    
</dd>
</dl>

<dl>
<dd>

**description:** `*string` — a description of the client
    
</dd>
</dl>

<dl>
<dd>

**additionalProperties:** `map[string]*generated.APIClientConfigRequestAdditionalPropertiesValue` — additional properties for the client
    
</dd>
</dl>

<dl>
<dd>

**ipFilterOp:** `*generated.APIClientConfigRequestIPFilterOp` — the operation of the ip filter. The default setting is null, which means that the ip filter is disabled
    
</dd>
</dl>

<dl>
<dd>

**ipFilters:** `[]string` — the list of ips
    
</dd>
</dl>

<dl>
<dd>

**jwkURI:** `*string` — the JSON web key URI endpoint
    
</dd>
</dl>

<dl>
<dd>

**additionalConfig:** `*generated.APIClientAdditionalConfig` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.APIClients.BulkDeleteAPIClient(request) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes all of the specified API clients. A successful response is returned even if the object does not exist.<br><br>Entitlements required: manageAPIClients (Manage API clients)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := []*generated.BulkOperation{
    &generated.BulkOperation{},
}
client.APIClients.BulkDeleteAPIClient(
    context.TODO(),
    request,
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**request:** `[]*generated.BulkOperation` 
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

