package appsignals

/* SIGNAL-LIST:
 * PAGE RELATED EVENTS
 *  onPageCreate(editor,pageTitle)
 *  onPageDelete(editor)
 *  onPageUpdate(editor)
 * USER RELATED EVENTS
 *  onCreateAccount(username)
*/

type SignalHandler func(interface{})
