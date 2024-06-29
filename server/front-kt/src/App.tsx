import './App.css';
import {Button, createTheme, Group, MantineProvider, SimpleGrid, Tabs, useMantineColorScheme} from '@mantine/core';
import '@mantine/core/styles.css';
import {HistoryPage} from "./pages/History";

const theme = createTheme({});

export default function App() {
    // <StompSessionProvider url={"ws://localhost:8080/ws"} /*debug={str => console.log(str)}*/>
    //
    //     </StompSessionProvider>
    return (
        <MantineProvider theme={theme} defaultColorScheme="dark">
                <Test/>
                {/*<TestSub/>*/}
        </MantineProvider>
    );
}

function Test() {
    const {setColorScheme, clearColorScheme} = useMantineColorScheme();

    return (
        <Tabs defaultValue="history">
            <Tabs.List>
                <Tabs.Tab value="history" leftSection="">
                    History
                </Tabs.Tab>
                <Tabs.Tab value="settings" leftSection="">
                    Settings
                </Tabs.Tab>
            </Tabs.List>

            <Tabs.Panel value="history">
                    <HistoryPage/>
            </Tabs.Panel>

            <Tabs.Panel value="settings">
                Settings tab content
                <Group>
                    <Button onClick={() => setColorScheme('light')}>Light</Button>
                    <Button onClick={() => setColorScheme('dark')}>Dark</Button>
                    <Button onClick={() => setColorScheme('auto')}>Auto</Button>
                    <Button onClick={clearColorScheme}>Clear</Button>
                </Group>
            </Tabs.Panel>
        </Tabs>
    )
}
