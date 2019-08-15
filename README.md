# magnificent

A new challenger appears!

## Installation

 1. Install Python. Python 2 and Python 3 both work.
 2. `pip install twisted`
 3. `python server.py`
 4. Okay, now you're running magnificent!
 5. Visit http://localhost:12345 in a web browser or something.
 6. It should throw a verbose error, or return "Magnificent!".

## Now it's your turn.

Requests to Magnificent usually succeed with a status code of 200, but we've noticed that requests can also fail about 25% of the time with a
status code of 500. The success and failure cases have fairly constrained return values and messages, but keep in mind that in a real production
system there could be any number of failure cases. We want to write a service to monitor the health of Magnificent.

Your user has tools to monitor log files, so that is one method you can use to get their attention.

### Rules

* This should go without saying, but you can't use anyone else's code, unless it's under [an open-source license](http://opensource.org/licenses).
  If you need to copy a snippet from Stack Overflow or some similar site, that's okay. But you should
  be the primary author of, and understand, all the code you submit.
* Please refrain from including any identifying information in your submission, such as comments with your name, author information in a setup file, etc.

### Requirements

* It should check Magnificent frequently - at least several times a minute.
* The service has to indicate how healthy the Magnificent service has been over the last little while, not only the success/failure of the most recent request.
* The service should also indicate if Magnificent has not been responding at all.


### Judging

This is what we look for.

* **Correctness**. Does it do what we asked?
* **Design**. Is the design appropriate for the task?
* **Clarity.** Can any competent programmer easily understand what's going on?
* **Generality.** It shouldn't be all hardcoded, but don't make it too abstract either.
* **Tests and testability.** It would be really great if you have tests. If not, it should be at least possible to test this.
* **Documentation:**
  * Can a junior developer get this running? Are requirements listed, including how to install them? What platform did you develop and test your solution on?
  * Is there *just enough* documentation, to tell us *why* the program works this way?

## Platform

Sauce Labs' services are usually written in Python and run on Linux.

However, please use the language and environment where you are the strongest.

**Note:** we would prefer that your solution be operating system agnostic. If it is not please mention so, and
be aware that your solution may take us slightly longer to review as we dig a Windows laptop out of the closet :)

For instance, you could use any of these:

    Python
    Go
    Shell scripts
    Ruby
    Node.JS


Please do not use any of these, as we do not have the expertise to
evaluate solutions in these environments:

    x86 Assembly
    Intercal
    Plan 9 from Bell Labs
    Coldfusion
    Malbolge
    Turbo Pascal
    NeXTSTEP
    Visual Basic for Applications
