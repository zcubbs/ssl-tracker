import {
    QueryClient,
    QueryClientProvider,
} from '@tanstack/react-query'
import {MantineProvider} from '@mantine/core';
import {Notifications} from '@mantine/notifications';

import DomainsTable from './DomainsTable';
import AddDomainForm from './AddDomainForm';

const queryClient = new QueryClient();

function App() {
    return (
        <QueryClientProvider client={queryClient}>
            <MantineProvider withNormalizeCSS withGlobalStyles>
                <Notifications/>
                <>
                    <h1>SSL Certificate Checker</h1>
                    <AddDomainForm/>
                    <DomainsTable/>
                </>
            </MantineProvider>
        </QueryClientProvider>
    );
}

export default App;

