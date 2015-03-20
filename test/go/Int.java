
public class Int extends LValue_int{

	public int value;
	private final int addr;

	public Int(int addr){
		this(addr, 0);
	}

	public Int(int addr, int value){
		this.addr = addr;
		this.value = value;
	}

	public int value(){
		return value;
	}

	public int addr(){
		return this.addr;
	}

	public void set(int value){
		this.value = value;
	}

	public boolean equals(Object other){
		return other instanceof Int 
			&& ((Int)other).value == this.value;	
	}
}
