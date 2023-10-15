import {
    QueryClient,
    QueryClientProvider,
} from '@tanstack/react-query'
import {Container, MantineProvider, Space, Text} from '@mantine/core';
import {Notifications} from '@mantine/notifications';

import DomainsTable from './DomainsTable';
import AddDomainForm from './AddDomainForm';

const queryClient = new QueryClient();

function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <MantineProvider
                withGlobalStyles
                withNormalizeCSS
                theme={{
                    colorScheme: 'dark',

                    headings: {
                        fontFamily: 'Roboto, sans-serif',
                        sizes: {
                            h1: { fontSize: '2rem' },
                        },
                    },
                }}
            >
                <Notifications/>
                <Container size="md" px="xs">
                    <h1>TLZ;</h1>
                    <Text>TLS Certificate Checker</Text>
                    <Space h="md" />
                    <AddDomainForm/>
                    <Space h="md" />
                    <DomainsTable/>
                </Container>
            </MantineProvider>
        </QueryClientProvider>
    );
}

export default App;

