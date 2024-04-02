import React, { useState } from 'react';
import useAxios from '../useHook/useHook'; // assuming you have defined the useAxios hook
import './Header.css';

const Header = () => {
    const [option, setOption] = useState('get');
    const [value1, setValue1] = useState('');
    const [value2, setValue2] = useState('');
    const [expiration, setExpiration] = useState('');


    const { response, error, isLoading, fetchData ,setResponse,setError,setIsLoading } = useAxios();

    const handleGetSubmit = () => {
        if (value1===""){
            setError("Key cannot be empty")
            return
        }

        setResponse(null)
        setError(null)
        setIsLoading(false)
        fetchData('GET', `http://localhost:8080/get/${value1}`);
    };

    const handleSetSubmit = () => {
        if (value1===""){
            setError("Key cannot be empty")
            return
        }

        if (value2===""){
            setError("value cannot be empty")
            return
        }

        setResponse(null)
        setError(null)
        setIsLoading(false)
        fetchData('PUT', 'http://localhost:8080/set', { "key": value1, "value": value2,"expiration": parseInt(expiration,10) });
    };

    return (
        <div className="header-container">
            <h1>DataHub</h1>
            <div className="header-options">
                <label>
                    <input
                        type="radio"
                        value="get"
                        checked={option === 'get'}
                        onChange={() => setOption('get')}
                    />
                    Get
                </label>
                <label>
                    <input
                        type="radio"
                        value="set"
                        checked={option === 'set'}
                        onChange={() => {setOption('set');setValue1('');setValue2('');setExpiration('')}}
                    />
                    Set
                </label>
            </div>
            <div className="header-form">
                {option === 'get' ? (
                    <div>
                        <input
                            type="text"
                            value={value1}
                            onChange={(e) => setValue1(e.target.value)}
                            placeholder="Enter Key"
                            required
                        />
                        <button onClick={handleGetSubmit}>Get Data</button>
                    </div>
                ) : (
                    <div>
                        <input
                            type="text"
                            value={value1}
                            onChange={(e) => setValue1(e.target.value)}
                            placeholder="Key"
                            required
                        />
                        <input
                            type="text"
                            value={value2}
                            onChange={(e) => setValue2(e.target.value)}
                            placeholder="Value"
                            required
                        />

                        <input
                            type="text"
                            value={expiration}
                            onChange={(e) => setExpiration(e.target.value)}
                            placeholder="Expiry"
                            required
                        />
                        <button onClick={handleSetSubmit}>Set Data</button>
                    </div>
                )}
            </div>
            <div className="header-response">
                {isLoading && <div>Loading...</div>}
                {error && <div>Error: {error}</div>}
                {response && <div> {JSON.stringify(response.data)}</div>}
            </div>
        </div>
    );
};

export default Header;