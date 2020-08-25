class Quotes extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        if (this.props.quotes) {
            return (
                <div class="text-center">
                    <h1>Latest</h1>
                    {this.props.quotes.map(quote => (
                        <div class="card text-center">
                            <div class="card-body">
                                <h5 class="card-title">{quote.symbol}</h5>
                                <p class="card-text" >${quote.price} <span style={{color: parseFloat(quote.changep) >= 0 ? "green" : "red"}}>({quote.changep})</span></p>
                            </div>
                        </div>
                    ))}
                </div>
            )
        } else {
            return (<div></div>)
        }
    }
}

class SingleStock extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            isLoaded: false,
            intra: []
        }
    }

    componentDidMount = () => {
        fetch('/api/intraday/' + this.props.stock)
            .then(res => res.json())
            .then((data) => {
                console.log("Data", data)
                this.setState({
                    isLoaded: true,
                    intra:data
                })
            })
            .catch(console.log)
    }

    render() {
        const { isLoaded, intra } = this.state

        if (!isLoaded) {
            return <div>Loading...</div>
        } else if (intra.prices.length < 1) {
            return (
                <div className="container">
                    <NavBar />
                    <div className="row text-center">
                        <center>
                            <h1>No data found for symbol '{this.props.stock.toUpperCase()}'</h1>
                            Make sure this is a valid symbol, or try your request again.
                        </center>
                    </div>
                </div>
            )
        } else {
            var lineColor
            if (intra.prices[0] < intra.prices[intra.prices.length - 1])
                lineColor = "green"
            else
                lineColor = "red"

            const opts = {
                title: this.props.stock.toUpperCase(),
                width: 400,
                height: 300,
                series: [
                    {},
                    {
                        stroke: lineColor
                    }
                ]
            };
            const data = [intra.times, intra.prices]

            return (
                <div className="container">
                    <NavBar />
                    <div className="row text-center">
                        <center>
                            <Chart options={opts} data={data} />
                        </center>
                    </div>
                </div>
            )
        }
    }
}

class Chart extends React.Component {
    constructor(props) {
        super(props)
    }

    componentDidMount() {
        this.u = new uPlot(this.props.options, this.props.data, this.el)
    }

    render() {
        return <div ref={el => (this.el = el)} />
    }
}

class NavBar extends React.Component {
    render() {
        return (
            <nav class="navbar navbar-default">
                <div class="container-fluid">
                    <div class="navbar-header">
                        <a class="navbar-brand" href="/">Stocks!</a>
                    </div>
                    <div class="collapse navbar-collapse" id="bs-example-navbar-collapse-2">
                        <ul class="nav navbar-nav navbar-right">
                            <li><a href="#">Invest with the wisdom of the crowd.</a></li>
                        </ul>
                    </div>
                </div>
            </nav>
        )
    }
}

class App extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            quotes: [{'symbol':'DEMO','price':50.00}],
            currentStock: null,
            value: '',
        }
        this.setStock = this.setStock.bind(this)
        this.handleSubmit = this.handleSubmit.bind(this)
        this.handleChange = this.handleChange.bind(this)
    }

    componentDidMount = () => {
        fetch('/api/list')
            .then(res => res.json())
            .then((data) => {
                console.log("Data", data)
                this.setState({ quotes:data })
            })
            .then(console.log("updated state",this.state))
            .catch(console.log)
    }

    setStock(symbol) {
        this.setState(state => ({ currentStock:symbol }))
    }

    handleChange(event) {
        this.setState({value: event.target.value});
    }

    handleSubmit(event) {
        console.log("submit:", this.state.value)
        this.setStock(this.state.value)
        event.preventDefault()
    }

    render () {
        if (this.state.currentStock !== null) {
            return <SingleStock stock={this.state.currentStock}/>
        }
        return (
            <div className="container-fluid">
                <NavBar />
                <div class="row text-center">
                    <div class="col-sm-4 col-sm-offset-4">
                        <form class="form-horizontal" onSubmit={this.handleSubmit} onChange={this.handleChange}>
                            <div class="form-group">
                                <div class="input-group">
                                    <input class="form-control" type="text" placeholder="Symbol (try TSLA)..." value={this.state.value} />
                                    <span class="input-group-btn">
                                        <button class="btn btn-primary" type="submit">Search</button>
                                    </span>
                                </div>
                            </div>
                        </form>
                    </div>
                </div>
                <Quotes quotes={this.state.quotes}/>
            </div>
        );
    }
}

ReactDOM.render(<App />, document.getElementById("app"));
