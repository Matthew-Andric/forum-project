# Portfolio Project - Forum
### Motivation/Reasoning for the Project
This project was something I made as I wanted to learn more about back end development and also to learn more about Go. The idea it started with was I was going to use Go and develop something that could interact with a database, I hadn't properly worked with a more modern language which was my motivation for picking Go and I just wanted to build more confidence with programming in general. Postgres was also used as the database in this case as I have some familiarity with it but wanted to use it with a project to learn more.

I was mostly focusing on the back end of this, front end is done in a basic way with mostly static HTML with some very light javascript.

### Demonstrations
User authentication/registration with conditional rendering based on a users role.
![](https://github.com/Matthew-Andric/forum-project/blob/main/examplegifs/login.gif)

Able to create/edit/delete posts and create/edit threads
![](https://github.com/Matthew-Andric/forum-project/blob/main/examplegifs/posteditdelete.gif)

Ability to upload and display profile pictures
![](https://github.com/Matthew-Andric/forum-project/blob/main/examplegifs/editprofilepicture.gif)

### Features Implemented
* User Authentication/Registration
  - Registration with credentials stored in the database
  - Bcrypt used for password salting/encryption
  - Sessions are tracked in the database with an encrypted string stored in a cookie with validation comparing the cookie to the stored valid sessions to confirm the user is allowed to perform only valid actions
* Dynamic Page Rendering
  - Boards/Threads/Posts/User Profiles are loaded from the database with the layout not being reliant on fixed links between pages
* Forum Features
  - Able to create/edit/delete posts
  - Able to create/edit threads
  - Able to create/edit/delete forum categories/subcategories and their properties via the admin panel
* User Profile Features
  - Able to upload and set a profile picture, image will be validated to be a valid .jpg/.png before being applied
  - Able to update password
