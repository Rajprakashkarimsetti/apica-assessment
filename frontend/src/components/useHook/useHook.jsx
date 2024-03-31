import { useState } from 'react';
import axios from 'axios';

const useAxios = () => {
  const [response, setResponse] = useState(null);
  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const fetchData = async (method, url, data = null, params = null) => {
    setIsLoading(true);
    try {
      let result;
      switch (method) {
        case 'GET':
          result = await axios.get(url, { params });
          break;
        case 'POST':
          result = await axios.post(url, data, { params });
          break;
        case 'PUT':
          result = await axios.put(url, data, { params });
          break;
        case 'DELETE':
          result = await axios.delete(url, { params });
          break;
        case 'PATCH':
          result = await axios.patch(url,  data);
          break;
        default:
          throw new Error('Invalid HTTP method');
      }
      console.log(result.data)
      setResponse(result.data);
    } catch (error) {
    
        setError(error.response.data.Reason); // Set error to the reason from the response data
        
    }
    setIsLoading(false);
  };

  return { response, error, isLoading, fetchData,setResponse,setError,setIsLoading };
};

export default useAxios;
