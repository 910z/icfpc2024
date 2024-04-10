import React, {useState} from 'react';

import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Alert from 'react-bootstrap/Alert';

import ButtonsShowcase from './showcases/Buttons';
import ToastsShowcase from './showcases/Toasts';
import {Nav, Navbar, NavDropdown} from "react-bootstrap";
import ThemeSwitch from "./components/ThemeSwitch";
import {HashRouter, Route, Routes} from "react-router-dom";
import Problems from "./pages/Problems";

function AlertDismissibleExample() {
    const [show, setShow] = useState(false);

    if (show) {
        return (
            <Alert variant="danger" onClose={() => setShow(false)} dismissible>
                <Alert.Heading>
                    I am an alert of type <span className="dangerText">danger</span>! But
                    my color is Teal!
                </Alert.Heading>
                <p>
                    By the way the button you just clicked is an{' '}
                    <span className="infoText">Info</span> button but is using the color
                    Tomato. Lorem ipsum dolor sit amet, consectetur adipisicing elit.
                    Accusantium debitis deleniti distinctio impedit officia reprehenderit
                    suscipit voluptatibus. Earum, nam necessitatibus!
                </p>
            </Alert>
        );
    }
    return (
        <Button variant="info" onClick={() => setShow(true)}>
            Show Custom Styled Alert
        </Button>
    );
}

const App = () => (
    <HashRouter>
        <div>
            <Navbar expand="lg" className="bg-body-tertiary">
                <Container>
                    <Navbar.Brand href={`/`}>ICFPC</Navbar.Brand>
                    <Navbar.Toggle aria-controls="basic-navbar-nav"/>
                    <Navbar.Collapse id="basic-navbar-nav">
                        <Nav className="me-auto">
                            <Nav.Link href={`#/problems`}>Problems</Nav.Link>
                            <Nav.Link href={`#/solutions`}>Solutions</Nav.Link>
                            <NavDropdown title="Dropdown" id="basic-nav-dropdown">
                                <NavDropdown.Item href="#action/3.1">Action</NavDropdown.Item>
                                <NavDropdown.Item href="#action/3.2">Another action</NavDropdown.Item>
                                <NavDropdown.Item href="#action/3.3">Something</NavDropdown.Item>
                                <NavDropdown.Divider/>
                                <NavDropdown.Item href="#action/3.4">Separated link</NavDropdown.Item>
                            </NavDropdown>
                        </Nav>
                    </Navbar.Collapse>
                    <ThemeSwitch></ThemeSwitch>
                </Container>
            </Navbar>
            <Container className="p-3">
                {/*<Container className="pb-1 p-5 mb-4 bg-light rounded-3">*/}
                {/*    <h1 className="header">Welcome To React-Bootstrap</h1>*/}
                {/*    <h2 className="header">Using Sass with custom theming</h2>*/}
                {/*    <AlertDismissibleExample />*/}
                {/*    <hr />*/}
                {/*    <p>*/}
                {/*        You can check further in information on the official Bootstrap docs{' '}*/}
                {/*        <a*/}
                {/*            href="https://getbootstrap.com/docs/4.3/getting-started/theming/#importing"*/}
                {/*            target="_blank"*/}
                {/*            rel="noopener noreferrer"*/}
                {/*        >*/}
                {/*            here*/}
                {/*        </a>*/}
                {/*        .*/}
                {/*    </p>*/}
                {/*</Container>*/}
                {/*<Container className="p-5 mb-4 bg-light rounded-3">*/}
                {/*    <h1 className="header">*/}
                {/*        Welcome To React-Bootstrap TypeScript Example*/}
                {/*    </h1>*/}
                {/*</Container>*/}
                <h2>Buttons</h2>
                <ButtonsShowcase/>
                <h2>Toasts</h2>
                <ToastsShowcase/>

                <Routes>
                    <Route path="/problems" element={<Problems/>}/>
                    <Route path="/solutions" element={<ButtonsShowcase/>}/>
                </Routes>
            </Container>
        </div>
    </HashRouter>
);

export default App;
