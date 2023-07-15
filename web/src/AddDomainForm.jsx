import React from 'react';
import axios from 'axios';
import {Button, TextInput} from '@mantine/core';
import { useForm } from '@mantine/form';

axios.defaults.baseURL = 'http://localhost:8000/api';

async function addDomain(data) {
    console.log(data);
    await axios.post('/domains', data);
}

function AddDomainForm() {
    const form = useForm({
        initialValues: {
            name: '',
        }
    });

    return (
        <form onSubmit={form.onSubmit((values) => addDomain(values))}>
            <TextInput
                label="Domain name"
                placeholder="Enter a domain to track"
                {...form.getInputProps('name')}
            />
            <Button type="submit">Track Domain</Button>
        </form>
    );
}

export default AddDomainForm;
