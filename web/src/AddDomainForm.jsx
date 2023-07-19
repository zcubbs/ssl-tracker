import axios from 'axios';
import {Button, Flex, TextInput} from '@mantine/core';
import {useForm} from '@mantine/form';
import {QueryClient, useQueryClient} from "@tanstack/react-query";

axios.defaults.baseURL = 'http://localhost:8000/api';

function AddDomainForm() {
    const form = useForm({
        initialValues: {
            name: '',
        }
    });
  const queryClient = useQueryClient();

  const addDomain = async (data) => {
    console.log(data);
    await axios.post('/domains', data).catch((err) => {
      console.log(err);
    }).then((res) => {
      console.log(res);
      form.reset();
      queryClient.refetchQueries(["domains"], { active: true })
    });
  }

    return (

        <form onSubmit={form.onSubmit((values) => addDomain(values))}>
            <Flex
                mih={50}
                gap="md"
                justify="flex-start"
                align="flex-start"
                direction="row"
                wrap="wrap"
            >
                <TextInput
                    placeholder="Enter a domain to track"
                    {...form.getInputProps('name')}
                />
                <Button type="submit">Track Domain</Button>
            </Flex>
        </form>
    );
}

export default AddDomainForm;
