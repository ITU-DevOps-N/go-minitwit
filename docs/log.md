# Session 02
## Why we chose Go + Gin Framework
We chose Google's relatively new programming language Golang for different reasons: 
- It is good for developing web applications. 
- It is easy to maintain.
- It is faster than many other programming languages and often compared to C in terms of speed and efficiency.
- Go does not require an interpreter, freeing up power and boosts performance.
- Gin is a high-performance Golang framework for creating web applications and contains a lot of useful modules.
- Gin reduces boilerplate code by providing commonly used functionalities (such as routing, middleware support, rendering, etc.), thus also making it easier to build web applications.
- Lastly, there was a personal wish of working with Golang as it is a new programming language that no one within the group has worked with before. 

## References
- https://yalantis.com/blog/why-use-go/
- https://github.com/gin-gonic/gin#gin-web-framework
---
# Session 03
## Why we chose DigitalOcean and Vagrant
*Remember to log and provide good arguments for the choice of virtualization techniques and deployment targets.*

- Coupled with GitHub education pack, DigitalOcean provides free credits, which is advantageous for us as students.
- DigitalOcean is a simple cloud service provider that is more suitable for small applications and projects, and simple to set up. 
- Combining DigitalOcean with Vagrant provides a clean way to set boot up virtual machines to host a web application.
- Documentation for DigitalOcean is straightforward and easily comprehensible to get started quickly.
- Vagrant is the bread and butter tool for interacting, setting up and managing virtual machines. Narrowed down it can be seen as a scripting engine for VirtualBox.
- Vagrant can be source controlled with ease, since everything is defined in a single text file. You can, of course, version control a snapshot of the virtual machine, but that wil take up a lot more space than just a Vagrantfile.
---
# Session 04
## Why we chose Travis CI

*Remember to log and provide good arguments for the choice of CI/CD system, i.e., why do you choose your solution instead of any other?*

- Coupled with GitHub education pack, Travis CI provides free credits, which is advantageous for us as students.
- Travis CI requires less setup than other alternatives (f.ex. Jenkins)
- Travis is a standalone solution for CI/CD, which can be beneficial if we e.g., decide to change cloud provider later. 
- Straightforward to configure through a YAML file

## Why we chose GORM
*Remember to log and provide good arguments for the choice of ORM framework and chosen DBMS*

- GORM is something we already integrated to our application when we chose to use Go as the primary language
- Able to automatically migrate schema - GORM automatically creates tables based on our model.
- GORM is developer friendly and well documented, which makes it easy to work with