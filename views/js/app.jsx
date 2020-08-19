class Quotes extends React.Component {
	constructor(props) {
		super(props)
	}

	render() {
		if (this.props.quotes) {
			return (
				<div>
					<center><h1>Latest</h1></center>
					{this.props.quotes.map(quote => (
						<div class="card text-center">
							<div class="card-body">
								<h5 class="card-title">{quote.symbol}</h5>
								{/*<h6 class="card-subtitle mb-2 text-muted">{quote.price}</h6>*/}
								<p class="card-text">{quote.price} ({quote.changep})</p>
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
		fetch('http://localhost:3000/api/intraday/' + this.props.stock)
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
	/*
	 *"2020-08-17 11:00:00": {
			"1. open": "125.3000",
			"2. high": "125.3750",
			"3. low": "125.2100",
			"4. close": "125.2200",
			"5. volume": "41169"
		},
		"2020-08-17 10:55:00": {
			"1. open": "125.0650",
			"2. high": "125.3400",
			"3. low": "125.0300",
			"4. close": "125.2900",
			"5. volume": "71923"
		},
		"2020-08-17 10:50:00": {
			"1. open": "124.9100",
			"2. high": "125.0600",
			"3. low": "124.9100",
			"4. close": "125.0300",
			"5. volume": "48198"
		},
		"2020-08-17 10:45:00": {
			"1. open": "124.7950",
			"2. high": "124.9500",
			"3. low": "124.7800",
			"4. close": "124.8900",
			"5. volume": "33126"
		},
		"2020-08-17 10:40:00": {
			"1. open": "124.9200",
			"2. high": "124.9600",
			"3. low": "124.7650",
			"4. close": "124.7750",
			"5. volume": "37698"
		},
		*/
	render() {
		const { isLoaded, intra } = this.state

		if (!isLoaded) {
			return <div>Loading...</div>
		} else {
			var lineColor
			if (intra.prices[0] < intra.prices[intra.prices.length - 1])
				lineColor = "green"
			else
				lineColor = "red"

			const opts = {
				title: this.props.stock,
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
					<div className="row">
						<div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
							<h1>Stonks!</h1>
							<p>Invest with the wisdom of the crowd.</p>
						</div>
					</div>
					<div className="row text-center">
						{/*<h1>Stonks! - {this.props.stock}</h1>*/}
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
		fetch('http://localhost:3000/api/list')
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
			<div className="container">
				<div className="row">
					<div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
						<h1>Stonks!</h1>
						<p>Invest with the wisdom of the crowd.</p>
					</div>
				</div>
				<div className="text-center">
					<form onSubmit={this.handleSubmit} onChange={this.handleChange}>
						<label>
							Stock: <input type="text" value={this.state.value} />
						</label>
						<input type="submit" value="Search" />
					</form>
				</div>
				<Quotes quotes={this.state.quotes}/>
			</div>
		);
	}
}

ReactDOM.render(<App />, document.getElementById("app"));
