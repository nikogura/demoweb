# DemoWeb

An unforgivably crude web app created to POC various web technologies.

# Pages
*/* Unauthenticated landing page

*/loggedin* Authenticated landing page

*/group-a* Page accessible only to members of group A

*/group-b* Page accessible only to members of group B

*/admin* Page accessible only to admins

# Users and Access

* anonymous users should be able to access `/`, but nothing else.

* userA should be able to access `/`, `/loggedin`, and `/group-a`.  `/group-b` and `/admin` should return 404

* userB should be able to access `/`, `/loggedin`, and `/group-b`.  `/group-a` and `/admin` should return 404

* admin should be able to access `/`, `/loggedin`, and `/group-a`, `/group-b` and `/admin`.
