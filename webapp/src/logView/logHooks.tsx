import { useEffect, useState } from "react";
import axios, { AxiosResponse } from "axios";


export interface LogQueryParams {
    minTime: string;
    maxTime: string;
    offset: number;
    limit: number;
    excludeLoggers: number[];
    search: string;
}

export interface LogEntry {
    id: number;
    timestamp: string;
    logger: string;
    loggerId: number;
    level: string;
    message: string;
    meta: any;
}

export interface LogFilters {
    get: LogQueryParams;
    set: {
        minTime: (time: string) => void;
        maxTime: (time: string) => void;
        offset: (offset: number) => void;
        limit: (limit: number) => void;
        excludeLoggers: {
            add: (loggerId: number) => void;
            remove: (loggerId: number) => void;
            raw: (loggerIds: number[]) => void;
        };
        search: (search: string) => void;
    };
    clear: () => void;
}


export const useFilters = () => {
    const [minTime, setMinTime] = useState<string>("");
    const [maxTime, setMaxTime] = useState<string>("");
    const [offset, setOffset] = useState<number>(0);
    const [limit, setLimit] = useState<number>(100);
    const [excludeLoggers, setExcludeLoggers] = useState<number[]>([]);
    const [search, setSearch] = useState<string>("");

    const addExcludeLogger = (loggerId: number) => {
        if (!excludeLoggers.includes(loggerId)) {
            setExcludeLoggers([...excludeLoggers, loggerId]);
        }
    }

    const removeExcludeLogger = (loggerId: number) => {
        setExcludeLoggers(excludeLoggers.filter((id) => id !== loggerId));
    }

    const logFilters = {
        get: {
            minTime,
            maxTime,
            offset,
            limit,
            excludeLoggers,
            search
        },
        set: {
            minTime: setMinTime,
            maxTime: setMaxTime,
            offset: setOffset,
            limit: setLimit,
            excludeLoggers: {
                add: addExcludeLogger,
                remove: removeExcludeLogger,
                raw: setExcludeLoggers
            },
            search: setSearch
        },
        clear: () => {
            setMinTime("");
            setMaxTime("");
            setOffset(0);
            setLimit(100);
            setExcludeLoggers([]);
            setSearch("");
        }
    }
    return logFilters;
}

export interface LogsApi {
    data: LogEntry[];
    loading: boolean;
    error: string | null;
    refetch: () => void;
}

export const useFetchLogs = (params: LogQueryParams): LogsApi => {
    const [logs, setLogs] = useState<LogEntry[]>([]);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [trigger, setTrigger] = useState<boolean>(false);
    useEffect(() => {
        const endpoint = `/logs`;
        setLoading(true);
        axios.get(endpoint, {
            params, paramsSerializer: (paramsObj) => {
                const searchParams = new URLSearchParams();
                for (const [key, value] of Object.entries(paramsObj)) {
                    if (Array.isArray(value)) {
                        value.forEach((item) => searchParams.append(key, item));
                    } else {
                        searchParams.append(key, value);
                    }
                }
                return searchParams.toString();
            }
        })
            .then((res: AxiosResponse) => {
                if (res.statusText != "OK") {
                    throw new Error(
                        "Logs request failed, status: " +
                        res.status +
                        " " +
                        res.statusText
                    );
                }
                return res.data;
            })
            .then((data) => {
                console.log('received log data', data);
                setLogs(data);
            })
            .catch((err) => {
                console.error('error fetching log data', err);
                setError(err.message);
            })
            .finally(() => {
                setLoading(false);
            });
    }, [trigger])
    return { data: logs, loading, error, refetch: () => setTrigger(!trigger) };
}

export const useFetchLogData = (logId: number | null) => {
    const [data, setData] = useState<any>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    useEffect(() => {
        if (logId) {
            setLoading(true);
            const endpoint = `/logs/${logId}`;
            axios.get(endpoint)
                .then((res) => {
                    if (res.statusText != "OK") {
                        throw new Error(
                            "Log request failed, status: " +
                            res.status +
                            " " +
                            res.statusText
                        );
                    }
                    return res.data;
                })
                .then((data) => {
                    console.log(`received data for log ${logId}`, data);
                    setData(data);
                })
                .catch((err) => {
                    setError(err.message);
                    console.error(`error fetching log data for log ${logId}`, err);
                })
                .finally(() => {
                    setLoading(false);
                });
        }
    }, [logId])
    return { data, loading, error };
}

