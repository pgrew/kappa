package ssh

import (
	"fmt"
	"reflect"
<<<<<<< Updated upstream
=======
	"strings"
>>>>>>> Stashed changes

	"github.com/subsilent/kappa/datamodel"
	"github.com/subsilent/kappa/skl"
)

// Session provides session and connection related information
type Session struct {
	namespace string
	user      datamodel.User
}

// Executor executes successfully parsed queries
type Executor struct {
	session  Session
	terminal Terminal
	system   datamodel.System
}

// Execute processes each statement
func (e *Executor) Execute(w *ResponseWriter, stmt skl.Statement) {

<<<<<<< Updated upstream
=======
	// Verify session has a user
	if e.session.user == nil {
		w.Fail(InternalServerError, "could not determine session user")
		return
	}

>>>>>>> Stashed changes
	switch stmt.NodeType() {
	case skl.UseNamespaceType:
		e.handleUseStatement(w, stmt)
	case skl.CreateNamespaceType:
		e.handleCreateNamespace(w, stmt)
	}
}

func (e *Executor) handleUseStatement(w *ResponseWriter, stmt skl.Statement) {
	use, ok := stmt.(*skl.UseStatement)
	if !ok {
<<<<<<< Updated upstream
		w.Fail(InvalidStatementType, "expected UseStatement, got %s instead", reflect.TypeOf(stmt))
=======
		w.Fail(InvalidStatementType, "expected *UseStatement, got %s instead", reflect.TypeOf(stmt))
>>>>>>> Stashed changes
		return
	}

	// Get user from session
	user := e.session.user
<<<<<<< Updated upstream
	if user == nil {
		w.Fail(InternalServerError, "could not access user data")
		return
	}
=======
>>>>>>> Stashed changes

	// Get namespace store
	namespaceStore, err := e.system.Namespaces()
	if err != nil {
		w.Fail(InternalServerError, "could not access namespace data")
		return
	}

	// Get namespace
	name := use.Namespace()

	// Verify namespace existence
	_, err = namespaceStore.Get(name)
	if err == datamodel.ErrNamespaceDoesNotExist {
		w.Fail(NamespaceDoesNotExist, name)
		return
	} else if err != nil {
		w.Fail(InternalServerError, "could not access namespace data")
		return
	}

	// If the user is an admin, grant access
	if user.IsAdmin() {
		e.session.namespace = name
		e.terminal.SetPrompt(fmt.Sprintf("kappa: %s> ", name))
		w.Success(OK, "")
		return
	}

	// Verify user has access to the namespace or is admin
	// 		If user has access, update session namespace and terminal
	// 		Otherwise, return access denied error
	for _, namespace := range user.Namespaces() {
		if namespace == name {
			e.session.namespace = name
			e.terminal.SetPrompt(fmt.Sprintf("kappa: %s> ", name))
			w.Success(OK, "")
			return
		}
	}

	// Otherwise, the user is not authorized
	w.Fail(Unauthorized, "")
}

// Only the admin can create root namespaces.
// Admin can also create sub-namespaces for any existing namespace.
// If the user is not the admin, they must have the 'create.namespace'
//  permission for the parent namespace.
// Root namespaces don't have any periods.
func (e *Executor) handleCreateNamespace(w *ResponseWriter, stmt skl.Statement) {
<<<<<<< Updated upstream
	create, ok := stmt.(*skl.CreateNamespaceStatement)
	if !ok {
		w.Fail(InvalidStatementType, "expected CreateNamespaceStatement, got %s instead", reflect.TypeOf(stmt))
		return
	}

	// Get user from session
	user := e.session.user
	if user == nil {
		w.Fail(InternalServerError, "could not access user data")
		return
	}

	// // Get namespace store
	// namespaceStore, err := e.system.Namespaces()
	// if err != nil {
	// 	w.Fail(InternalServerError, "could not access namespace data")
	// 	return
	// }

	w.Success(OK, create.Namespace())

	// // Get namespace
	// name := create.Namespace()

	// // Verify namespace existence
	// _, err := namespaceStore.Get(name)

	// // If err == nil, the namespace already existed
	// if err == nil {
	// 	w.Fail(NamespaceAlreadyExists, name)
	// 	return
	// }

	//

	// // Create new namespace
	// namespace, err := namespaceStore.Create(name)

	// // If err !+ nil, namespace could not be created
	// if err == datamodel.ErrNamespaceDoesNotExist {
	// 	w.Fail(NamespaceDoesNotExist, name)
	// 	return
	// } else if err != nil {
	// 	w.Fail(InternalServerError, "could not access namespace data")
	// 	return
	// }

	// // If the user is an admin, grant access
	// if user.IsAdmin() {
	// 	e.session.namespace = name
	// 	e.terminal.SetPrompt(fmt.Sprintf("kappa: %s> ", name))
	// 	w.Success(OK, "")
	// 	return
	// }

	// // Verify user has access to the namespace or is admin
	// // 		If user has access, update session namespace and terminal
	// // 		Otherwise, return access denied error
	// for _, namespace := range user.Namespaces() {
	// 	if namespace == name {
	// 		e.session.namespace = name
	// 		e.terminal.SetPrompt(fmt.Sprintf("kappa: %s> ", name))
	// 		w.Success(OK, "")
	// 		return
	// 	}
	// }

	// // Otherwise, the user is not authorized
	// w.Fail(Unauthorized, "")
=======

	createStatement, ok := stmt.(*skl.CreateNamespaceStatement)
	if !ok {
		w.Fail(InvalidStatementType, "expected *CreateNamespaceStatement, got %s instead", reflect.TypeOf(stmt))
		return
	}

	// Get namespace store
	namespaceStore, err := e.system.Namespaces()
	if err != nil {
		w.Fail(InternalServerError, "could not access namespace data")
		return
	}

	// Get session user
	user := e.session.user

	// Get namespace
	namespace := createStatement.Namespace()

	// If err == nil, the namespace already existed
	if e.namespaceAlreadyExists(namespace, namespaceStore) {
		w.Success(NamespaceAlreadyExists, namespace)
		return
	}

	// If root namespace
	if createStatement.IsRootNamespace() {
		e.handleCreateRootNamespace(w, createStatement, namespaceStore)
		return
	}

	// Admin user is granted access by default
	access := user.IsAdmin()

	// Get parent namespace
	var parent datamodel.Namespace
	index := strings.LastIndex(namespace, ".")
	parentNamespace := namespace[:index]

	// If the user is not an admin check their permissions for the parent namespace
	if !access {

		// Get user roles for parent namespace
		roles := user.Roles(parentNamespace)

		// Determine if parent namespace exists
		ns, err := namespaceStore.Get(parentNamespace)
		if err == datamodel.ErrNamespaceDoesNotExist {
			w.Fail(NamespaceDoesNotExist, parentNamespace)
			return
		} else if err != nil {
			w.Fail(InternalServerError, "")
			return
		}

		// Memoize parent namespace
		parent = ns

		// Scan roles for permissions
		for _, role := range roles {
			if ns.HasPermission(role, createStatement.RequiredPermissions()) {
				access = true
			}
		}

		// Return error if not authorized
		if !access {
			w.Fail(Unauthorized, "cannot create namespace '%s'", namespace)
			return
		}
	}

	// If we've gotten this far, the user has permission to create the namespace

	// Get parent namespace
	if parent == nil {

		// Verify namespace existance
		parent, err = namespaceStore.Get(parentNamespace)
		if err != nil {
			w.Fail(InternalServerError, "parent namespace does not exist")
			return
		}
	}

	// Create child namespace
	if _, err = parent.CreateChild(namespace); err != nil {
		w.Fail(CreateNamespaceError, "cannot create namespace '%s'", namespace)
		return
	}

	w.Success(OK, "namespace created")
}

// namespaceAlreadyExists determines if a namespace already exists...
func (e *Executor) namespaceAlreadyExists(namespace string, store datamodel.NamespaceStore) bool {
	_, err := store.Get(namespace)
	return err == nil
}

// If the namespace being created is a root namespace, only the admin account can create it
func (e *Executor) handleCreateRootNamespace(w *ResponseWriter, stmt *skl.CreateNamespaceStatement, store datamodel.NamespaceStore) {

	// Get namespace
	name := stmt.Namespace()

	// Verify namespace existance
	_, err := store.Get(name)

	// If err == nil, the namespace already exists
	if err == nil {
		w.Success(NamespaceAlreadyExists, name)
		return
	}

	// Get session user
	user := e.session.user
	if user.IsAdmin() {

		// Create new namespace
		_, err := store.Create(name)

		// If err !+ nil, namespace could not be created
		if err != nil {
			w.Fail(CreateNamespaceError, "could not create namespace '%s'", name)
			return
		}

		// No error == success
		w.Success(OK, "namespace created")
		return
	}

	// Otherwise fail creation
	w.Fail(Unauthorized, "root namespaces can only be created by the admin account")
	return
>>>>>>> Stashed changes
}
